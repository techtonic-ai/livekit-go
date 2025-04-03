package main

import (
	"log"
	"os"

	"livekitgo/livekit"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	ROOM_ID = "demo-room"
)

func main() {
	app := fiber.New()

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://livekit-server-ydpb.onrender.com,http://localhost:3000,http://localhost:8000,http://192.168.1.9:8000",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Routes
	app.Get("/api/token/:participant", getToken)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Serve static files
	app.Static("/", "./")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Log environment information
	log.Printf("Server starting on port %s", port)
	log.Printf("Environment: %s", os.Getenv("ENVIRONMENT"))
	log.Printf("Render.com URL: %s", os.Getenv("RENDER_EXTERNAL_URL"))
	log.Printf("LiveKit WebSocket URL: %s", os.Getenv("LIVEKIT_WS_URL"))
	log.Fatal(app.Listen(":" + port))
}

func getToken(c *fiber.Ctx) error {
	participant := c.Params("participant")
	if participant != "participant1" && participant != "participant2" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid participant. Must be 'participant1' or 'participant2'",
		})
	}

	// Generate LiveKit token
	token, err := livekit.GenerateToken(ROOM_ID, participant)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Determine WebSocket URL based on environment
	wsHost := os.Getenv("LIVEKIT_WS_URL")
	if wsHost == "" {
		// Check if we're running on Render.com
		if os.Getenv("RENDER_EXTERNAL_URL") != "" {
			wsHost = "wss://livekit-server-ydpb.onrender.com"
		} else {
			// Default to local development
			wsHost = "ws://192.168.1.9:7881"
		}
	}

	log.Printf("\n=== LiveKit Connection Details ===")
	log.Printf("Room ID: %s\n", ROOM_ID)
	log.Printf("Participant: %s\n", participant)
	log.Printf("Token: %s\n", token)
	log.Printf("WebSocket URL: %s\n", wsHost)
	log.Printf("Environment: %s\n", os.Getenv("ENVIRONMENT"))
	log.Printf("Render.com URL: %s\n", os.Getenv("RENDER_EXTERNAL_URL"))
	log.Printf("===============================\n")

	return c.JSON(fiber.Map{
		"room_id": ROOM_ID,
		"token":   token,
		"ws_url":  wsHost,
	})
}
