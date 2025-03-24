package redis

import (
	"context"
	"fmt"
	"sync"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/internal"
	"github.com/redis/go-redis/v9"
)

var (
	mu          = &sync.Mutex{}
	once        sync.Once
	redisClient *redis.Client
)

func ConnectToRedis(ctx context.Context) (*redis.Client, error) {
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		})
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error during pinging Redis")
	}
	return redisClient, nil
}

func SetValue(ctx context.Context, key, value string) error {
	mu.Lock()
	defer mu.Unlock()

	if err := redisClient.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set %s-%s", key, value)
	}
	return nil
}

func GetValue(ctx context.Context, key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	fetched_value, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s from redis", key)
	}

	return fetched_value, nil
}

func AddOrUpdatePlayerQueueMMR(ctx context.Context, player *models.Player) error {
	mu.Lock()
	defer mu.Unlock()

	err := redisClient.ZAdd(ctx, player.QueueID, redis.Z{
		Score:  player.MatchmakingRating,
		Member: player.ID,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add player to the QueueID-%s", player.QueueID)
	}
	return nil
}

func RemovePlayerFromQueueMMR(ctx context.Context, player *models.Player) error {
	mu.Lock()
	defer mu.Unlock()

	if err := redisClient.ZRem(ctx, player.QueueID, player.ID).Err(); err != nil {
		return fmt.Errorf("failed to remove user %s from queue %s", player.ID, player.QueueID)
	}
	return nil
}

func GetAllQueueMemeber(ctx context.Context, queue_id string) ([]models.AbstractPlayer, error) {
	mu.Lock()
	defer mu.Unlock()

	members, err := redisClient.ZRangeWithScores(ctx, queue_id, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get memebers of queue-%s", queue_id)
	}
	var players []models.AbstractPlayer
	for _, member := range members {
		players = append(players, models.AbstractPlayer{
			ID:    member.Member.(string),
			Score: member.Score,
		})
	}
	return players, nil
}
