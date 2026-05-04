package store

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rishabh21g/magic-board/internal/domain"
)

// create a RedisStore struct that implements the Store interface
type RedisStore struct {
	client *redis.Client
}

// constructor function for creating a new RedisStore instance
func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) ClaimBlock(ctx context.Context, blockID, userID string) (bool, error) {
	ok, err := r.client.HSetNX(ctx, "grid", blockID, userID).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (r *RedisStore) GetOwner(ctx context.Context, blockID string) (string, error) {
	owner, err := r.client.HGet(ctx, "grid", blockID).Result()
	if err == redis.Nil {
		return "", nil // unclaimed
	}
	if err != nil {
		return "", err
	}
	return owner, nil
}

// get entire grid
func (r *RedisStore) GetAllBlocks(ctx context.Context) ([]*domain.Block, error) {
	data, err := r.client.HGetAll(ctx, "grid").Result()

	if err != nil {
		return nil, errors.New(err.Error())
	}
	blocks := make([]*domain.Block, 0)
	for blockID, ownerID := range data {
		blocks = append(blocks, &domain.Block{
			BlockID:   blockID,
			OwnerID:   ownerID,
			Timestamp: 0,
		})
	}
	return blocks, nil
}

func (r *RedisStore) UnclaimBlock(ctx context.Context, blockID string) (bool, error) {
	r.client.HDel(ctx, "grid", blockID).Result()
	return true, nil

}

func (r *RedisStore) GetLeaderBoard(ctx context.Context) ([]*domain.LeaderboardEntry, error) {
	data, err := r.client.HGetAll(ctx, "grid").Result()
	if err != nil {
		return nil, errors.New(err.Error())
	}
	counts := make(map[string]int)
	for _, ownerID := range data {
		counts[ownerID]++
	}
	leaderboard := make([]*domain.LeaderboardEntry, 0)
	for ownerID, count := range counts {
		leaderboard = append(leaderboard, &domain.LeaderboardEntry{
			OwnerID: ownerID,
			Count:   count,
		})
	}
	return leaderboard, nil
}
