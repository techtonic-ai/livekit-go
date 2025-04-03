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
		AllowOrigins:     "*",
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

	log.Printf("Server starting on port %s", port)
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
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	wsHost := os.Getenv("LIVEKIT_WS_URL")
	if wsHost == "" {
		wsHost = "wss://livekit-server-ydpb.onrender.com"
	}

	log.Printf("\n=== LiveKit Connection Details ===")
	log.Printf("Room ID: %s\n", ROOM_ID)
	log.Printf("Participant: %s\n", participant)
	log.Printf("Token: %s\n", token)
	log.Printf("WebSocket URL: %s\n", wsHost)
	log.Printf("===============================\n")

	return c.JSON(fiber.Map{
		"room_id": ROOM_ID,
		"token":   token,
		"ws_url":  wsHost,
	})
}
