package redis

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/redis/go-redis/v9"
)

var (
	mu          sync.Mutex
	once        sync.Once
	redisClient *redis.Client
)

type RedisClient struct {
	client *redis.Client
}

func ConnectToRedis(ctx context.Context) (*RedisClient, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	var initErr error
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		})

		if client == nil {
			initErr = fmt.Errorf("failed to create Redis client")
			return
		}

		_, pingErr := client.Ping(ctx).Result()
		if pingErr != nil {
			initErr = fmt.Errorf("failed to ping Redis: %w", pingErr)
			return
		}

		redisClient = client
	})

	if initErr != nil {
		return nil, initErr
	}

	if redisClient == nil {
		return nil, fmt.Errorf("redis client is nil after initialization")
	}

	return &RedisClient{client: redisClient}, nil
}

func (r *RedisClient) SetValue(ctx context.Context, key, value string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	if err := r.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("failed to set %s-%s: %w", key, value, err)
	}
	return nil
}

func (r *RedisClient) GetValue(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("redis client is nil")
	}

	mu.Lock()
	defer mu.Unlock()

	fetched_value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s from redis: %w", key, err)
	}

	return fetched_value, nil
}

func (r *RedisClient) AddOrUpdatePlayerQueueMMR(ctx context.Context, player *models.Player) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	if player == nil {
		return fmt.Errorf("player is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	exist, _ := r.client.ZScore(ctx, player.QueueID, player.Username).Result()
	if exist != 0 {
		return fmt.Errorf("player already exists in the queue with score: %f", exist)
	}

	err := r.client.ZAdd(ctx, player.QueueID, redis.Z{
		Score:  player.MatchmakingRating,
		Member: player.Username,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to add player to the QueueID-%s: %w", player.QueueID, err)
	}
	return nil
}

func (r *RedisClient) RemovePlayerFromQueueMMR(ctx context.Context, player *models.Player) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	if player == nil {
		return fmt.Errorf("player is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	exist, _ := r.client.ZScore(ctx, player.QueueID, player.Username).Result()
	if exist == 0 {
		return fmt.Errorf("player does not exist in the queue with score: %f", exist)
	}

	if err := r.client.ZRem(ctx, player.QueueID, player.Username).Err(); err != nil {
		return fmt.Errorf("failed to remove user %s from queue %s: %w", player.Username, player.QueueID, err)
	}
	return nil
}

func (r *RedisClient) RemoveBatchPlayersFromQueueMMR(ctx context.Context, players []models.AbstractPlayer) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	if players == nil {
		return fmt.Errorf("players is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	// Get the queue ID from the first player
	if len(players) == 0 {
		return nil
	}
	queueID := players[0].QueueID

	for _, player := range players {
		if err := r.client.ZRem(ctx, queueID, player.ID).Err(); err != nil {
			fmt.Printf("failed to remove user %s from queue %s: %v\n", player.ID, queueID, err)
			continue
		}
	}
	return nil
}

func (r *RedisClient) RemovePlayerByUsername(ctx context.Context, username string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	exist, _ := r.client.ZScore(ctx, "queue_1", username).Result()
	if exist == 1 {
		fmt.Printf("player does not exist in the queue with score: %f\n", exist)
	}
	if err := r.client.ZRem(ctx, "queue_1", username).Err(); err != nil {
		fmt.Printf("failed to remove user %s from queue %s: %v\n", username, "queue_1", err)
	}
	fmt.Printf("User: %s removed from queue\n", username)
	return nil
}

func (r *RedisClient) GetAllQueueMembers(ctx context.Context, queue_id string) ([]models.AbstractPlayer, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	members, err := r.client.ZRangeWithScores(ctx, queue_id, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get memebers of queue-%s: %w", queue_id, err)
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

func (r *RedisClient) GetAllQueueMembersLength(ctx context.Context, queue_id string) (int64, error) {
	if r.client == nil {
		return 0, fmt.Errorf("redis client is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	members, err := r.client.ZRange(ctx, queue_id, 0, -1).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get memebers of queue-%s: %w", queue_id, err)
	}
	return int64(len(members)), nil
}

func (r *RedisClient) GetPlayersByRatingRange(ctx context.Context, queue_id string, start, end int64, count_threshold int64) ([]models.AbstractPlayer, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redis client is nil")
	}
	mu.Lock()
	defer mu.Unlock()

	members, err := r.client.ZRangeByScoreWithScores(ctx, queue_id, &redis.ZRangeBy{
		Min:    strconv.FormatInt(start, 10),
		Max:    strconv.FormatInt(end, 10),
		Offset: 0,
		Count:  count_threshold,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get memebers of queue-%s: %w", queue_id, err)
	}

	var players []models.AbstractPlayer
	for _, member := range members {
		players = append(players, models.AbstractPlayer{
			ID:      member.Member.(string),
			Score:   member.Score,
			QueueID: queue_id,
		})
	}
	return players, nil
}
