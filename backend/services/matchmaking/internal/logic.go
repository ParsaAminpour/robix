package internal

import (
	"context"
	"sync"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/fatih/color"
)

// goroutine No.1 -> rating from 0 - 200
// goroutine No.2 -> rating from 200 - 400
// goroutine No.3 -> rating from 400 - 600
// goroutine No.4 -> rating from 600 - 800
// goroutine No.5 -> rating from 800 - 1000
func MatchMakeByRatingRange(wg *sync.WaitGroup, channel chan<- []models.AbstractPlayer, ctx context.Context, redis_client *redis.RedisClient, queue_id string, start, end int64, count_threshold int64) error {
	players, err := redis_client.GetPlayersByRatingRange(ctx, queue_id, start, end, count_threshold)
	if err != nil {
		return err
	}

	color.Green("players found in queue %s: %v", queue_id, players)
	channel <- players
	wg.Done()
	return nil
}
