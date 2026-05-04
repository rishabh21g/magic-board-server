package store

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

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
		return "", nil
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
			BlockID: blockID,
			OwnerID: ownerID,
		})
	}
	return blocks, nil
}

// unclaim a block by deleting the corresponding field from the Redis hash
func (r *RedisStore) UnclaimBlock(ctx context.Context, blockID string) (bool, error) {
	n, err := r.client.HDel(ctx, "grid", blockID).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// get top5 leaderboard entries by counting the number of blocks claimed by each user and returning the top 5 users with the most blocks claimed, sorted by count in descending order and then by userID in ascending order
func (r *RedisStore) GetLeaderBoard(ctx context.Context) ([]*domain.LeaderboardEntry, error) {
	data, err := r.client.HGetAll(ctx, "grid").Result()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	counts := make(map[string]int)
	for _, ownerID := range data {
		if ownerID == "" {
			continue
		}
		counts[ownerID]++
	}

	leaderboard := make([]*domain.LeaderboardEntry, 0, len(counts))
	for ownerID, count := range counts {
		leaderboard = append(leaderboard, &domain.LeaderboardEntry{
			OwnerID: ownerID,
			Count:   count,
		})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].Count != leaderboard[j].Count {
			return leaderboard[i].Count > leaderboard[j].Count
		}
		return leaderboard[i].OwnerID < leaderboard[j].OwnerID
	})

	if len(leaderboard) > 5 {
		leaderboard = leaderboard[:5]
	}

	ids := make([]string, 0, len(leaderboard))
	for _, e := range leaderboard {
		ids = append(ids, e.OwnerID)
	}

	profiles, err := r.GetUserProfiles(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, e := range leaderboard {
		if p, ok := profiles[e.OwnerID]; ok && p != nil {
			e.Username = p.Username
			e.Color = p.Color
		}
	}

	return leaderboard, nil
}

// upsert user profile by storing the profile data as a JSON string in a Redis hash with the userID as the field and the JSON string as the value
func (r *RedisStore) UpsertUserProfile(ctx context.Context, userID, username, color string) error {
	if userID == "" {
		return errors.New("userID is required")
	}

	profile := domain.UserProfile{
		Username: username,
		Color:    color,
	}

	b, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return r.client.HSet(ctx, "users", userID, string(b)).Err()
}

// get user profiles by retrieving the JSON strings from the Redis hash for the given userIDs, unmarshaling them into UserProfile structs, and returning a map of userID to UserProfile
func (r *RedisStore) GetUserProfiles(ctx context.Context, userIDs []string) (map[string]*domain.UserProfile, error) {
	profiles := make(map[string]*domain.UserProfile)

	if len(userIDs) == 0 {
		return profiles, nil
	}

	unique := make([]string, 0, len(userIDs))
	seen := make(map[string]struct{}, len(userIDs))
	for _, id := range userIDs {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	if len(unique) == 0 {
		return profiles, nil
	}

	raw, err := r.client.HMGet(ctx, "users", unique...).Result()
	if err != nil {
		return nil, err
	}

	for i, v := range raw {
		if v == nil {
			continue
		}

		var s string
		switch t := v.(type) {
		case string:
			s = t
		case []byte:
			s = string(t)
		default:
			continue
		}

		var p domain.UserProfile
		if err := json.Unmarshal([]byte(s), &p); err != nil {
			return nil, err
		}

		profiles[unique[i]] = &p
	}

	return profiles, nil
}
