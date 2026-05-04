package game

import (
	"context"
	"fmt"
	"time"

	"github.com/rishabh21g/magic-board/internal/domain"
	"github.com/rishabh21g/magic-board/internal/store"
)

type Service struct {
	store store.Store
}

func NewService(store store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) ClaimBlock(ctx context.Context, blockID, userID string) (*domain.Block, error) {
	success, err := s.store.ClaimBlock(ctx, blockID, userID)
	fmt.Printf("Attempting to claim block %s for user %s: success=%v, error=%v\n", blockID, userID, success, err)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, fmt.Errorf("block %s is already claimed", blockID)
	}
	block := &domain.Block{
		BlockID:   blockID,
		OwnerID:   userID,
		Timestamp: time.Now().Unix(),
	}
	return block, nil
}
