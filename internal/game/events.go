package game

type ClaimBlockEvent struct {
	BlockID  string `json:"blockID"`
	UserID   string `json:"userID"`
	Username string `json:"username"`
	Color    string `json:"color"`
}

type UnclaimBlockEvent struct {
	BlockID string `json:"blockID"`
	UserID  string `json:"userID"`
}
