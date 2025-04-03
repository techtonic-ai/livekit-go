package livekit

import (
	"fmt"
	"time"

	"github.com/livekit/protocol/auth"
)

const (
	apiKey    = "devkey"
	apiSecret = "devsecret123456789012345678901234567890"
)

// GenerateToken creates a new token for a participant
func GenerateToken(roomID, participantID string) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSecret).
		AddGrant(&auth.VideoGrant{
			RoomJoin: true,
			Room:     roomID,
		}).
		SetIdentity(participantID).
		SetValidFor(24 * time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}
