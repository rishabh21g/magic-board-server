package domain

import "encoding/json"

type Block struct {
	BlockID   string `json:"blockID"`
	OwnerID   string `json:"userID"`
	Timestamp int64  `json:"timestamp"`
}

type LeaderboardEntry struct {
	OwnerID string `json:"userID"`
	Count   int    `json:"count"`
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
