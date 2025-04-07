package internal

import (
	"context"
	"sync"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
)

// goroutine No.1 -> rating from 0 - 200
// goroutine No.2 -> rating from 200 - 400
// goroutine No.3 -> rating from 400 - 600
// goroutine No.4 -> rating from 600 - 800
// goroutine No.5 -> rating from 800 - 1000
func MatchMakeByRatingRange(wg *sync.WaitGroup, channel chan<- []models.AbstractPlayer, ctx context.Context, redis_client *redis.RedisClient) error {
	players, err := redis_client.GetPlayersByRatingRange(ctx, "queue_mmr", 0, 200, 5)
	if err != nil {
		return err
	}
	channel <- players
	wg.Done()
	return nil
}
