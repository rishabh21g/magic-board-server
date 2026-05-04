package game

type ClaimBlockEvent struct {
	BlockID string `json:"blockID"`
	UserID  string `json:"userID"`
}

type UnclaimBlockEvent struct {
	BlockID string `json:"blockID"`
	UserID  string `json:"userID"`
}
