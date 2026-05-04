package domain

import (
	"encoding/json"
)

type Block struct {
	BlockID string `json:"blockID"`
	OwnerID string `json:"userID"`
}

type LeaderboardEntry struct {
	OwnerID string `json:"userID"`
	Count   int    `json:"count"`
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type UserProfile struct {
	Username string `json:"username"`
	Color    string `json:"color"`
}
