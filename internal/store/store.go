package store

import (
	"context"

	"github.com/rishabh21g/magic-board/internal/domain"
)

type Store interface {
	ClaimBlock(ctx context.Context, blockID, userID string) (bool, error)
	GetOwner(ctx context.Context, blockID string) (string, error)
	GetAllBlocks(ctx context.Context) ([]*domain.Block, error)
	UnclaimBlock(ctx context.Context, blockID string) (bool, error)
	GetLeaderBoard(ctx context.Context) ([]*domain.LeaderboardEntry, error)
}
