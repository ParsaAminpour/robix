package tests

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var redis_client *redis.RedisClient
var queue_id string = uuid.NewString()
var addedPlayerUntilNow int = 0
var wg sync.WaitGroup

func TestRedisConnectionAndPing(t *testing.T) {
	redis_client, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	// Test ping through SetValue and GetValue methods
	err = redis_client.SetValue(ctx, "test", "ping")
	assert.NoError(t, err)
	value, err := redis_client.GetValue(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, value, "ping")
}

func TestRedisSetAndGetKeyValue(t *testing.T) {
	redis_client, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	err = redis_client.SetValue(ctx, "name", "john")
	assert.NoError(t, err)

	fetched_value, fetched_err := redis_client.GetValue(ctx, "name")
	assert.NoError(t, fetched_err)
	assert.Equal(t, fetched_value, "john")
}

func TestRedisZAdd(t *testing.T) {
	redis_client, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	player := &models.Player{}
	another_player := &models.Player{}
	player = player.NewPlayer("john", queue_id, 10)
	another_player = another_player.NewPlayer("Elliot", queue_id, 20)

	t.Run("Test ZAdd", func(t *testing.T) {
		err = redis_client.AddOrUpdatePlayerQueueMMR(ctx, player)
		assert.NoError(t, err)
		addedPlayerUntilNow++
	})

	t.Run("Test ZRem", func(t *testing.T) {
		err = redis_client.AddOrUpdatePlayerQueueMMR(ctx, another_player)
		assert.NoError(t, err)
		addedPlayerUntilNow++

		var abs_players []models.AbstractPlayer
		abs_players, _ = redis_client.GetAllQueueMembers(ctx, queue_id)
		assert.Equal(t, len(abs_players), addedPlayerUntilNow) // one fron previous test and one from this test

		err = redis_client.RemovePlayerFromQueueMMR(ctx, player)
		assert.NoError(t, err)
		addedPlayerUntilNow--

		players, err := redis_client.GetAllQueueMembers(ctx, player.QueueID)
		assert.NoError(t, err)
		assert.Equal(t, players[0].ID, another_player.Username)
		assert.Equal(t, players[0].Score, another_player.MatchmakingRating)
		assert.Len(t, players, 1) // only one player is in queue
	})
}

func TestAddBatchPlayer(t *testing.T) {
	t.Run("add batch players to the queue", func(t *testing.T) {
		redis_client, err := redis.ConnectToRedis(ctx)
		assert.NoError(t, err)

		var abs_players []models.AbstractPlayer
		abs_players, _ = redis_client.GetAllQueueMembers(ctx, queue_id)
		for _, p := range abs_players {
			fmt.Printf("Player ID: %s\n", p.ID)
			fmt.Printf("Player score: %v\n", p.Score)
		}

		file, err := os.Open("./mocksData/MOCK_DATA.csv")
		assert.NoError(t, err)
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		assert.NoError(t, err)

		var players []models.Player
		for _, record := range records {
			username := record[0]
			score, _ := strconv.ParseFloat(record[1], 64)

			player := models.Player{
				Username:          username,
				MatchmakingRating: score,
				QueueID:           queue_id,
			}
			players = append(players, player)
		}

		startTime := time.Now()
		for _, player := range players {
			redis_client.AddOrUpdatePlayerQueueMMR(ctx, &player)
			// assert.NoError(t, err)
		}
		endTime := time.Now()
		addedPlayerUntilNow += len(players)
		color.Green("Time taken to add %d players in Redis: %v\n", len(players), endTime.Sub(startTime))

		abs_players, err = redis_client.GetAllQueueMembers(ctx, queue_id)
		assert.NoError(t, err)
		assert.Equal(t, len(abs_players), addedPlayerUntilNow)
	})

	t.Run("remove batch players from the queue", func(t *testing.T) {
		redis_client, err := redis.ConnectToRedis(ctx)
		assert.NoError(t, err)

		ch_players := make(chan []models.AbstractPlayer, 100)
		ended_channel := make(chan bool, 10)
		var endedOnce sync.Once
		var totalRemoved int64
		var mu sync.Mutex

		catchByRange := func(start, end int64) {
			defer wg.Done()
			emptyCount := 0
		CATCHER_LOOP:
			for {
				select {
				case <-ended_channel:
					break CATCHER_LOOP
				default:
					players, err := redis_client.GetPlayersByRatingRange(ctx, queue_id, start, end, 50)
					if err != nil {
						log.Printf("Error getting players: %v", err)
						continue
					}
					if len(players) > 0 {
						select {
						case ch_players <- players:
							mu.Lock()
							totalRemoved += int64(len(players))
							mu.Unlock()
							emptyCount = 0 // Reset empty count when we find players
						case <-ended_channel:
							break CATCHER_LOOP
						}
					} else {
						emptyCount++
						if emptyCount >= 3 { // If we get empty results 3 times in a row, consider this range done
							endedOnce.Do(func() {
								ended_channel <- true
							})
							break CATCHER_LOOP
						}
					}
					time.Sleep(100 * time.Millisecond) // Add a small delay to prevent tight loop
				}
			}
		}

		removeCatched := func() {
			defer wg.Done()
		REMOVE_LOOP:
			for {
				select {
				case players := <-ch_players:
					if err := redis_client.RemoveBatchPlayersFromQueueMMR(ctx, players); err != nil {
						log.Printf("Error removing players: %v", err)
					}
				case <-ended_channel:
					break REMOVE_LOOP
				}
			}
		}

		// Get initial count
		initialCount, err := redis_client.GetAllQueueMembersLength(ctx, queue_id)
		assert.NoError(t, err)
		fmt.Printf("Initial player count: %d\n", initialCount)

		startTime := time.Now()
		// Start with smaller ranges to ensure we don't miss any players
		wg.Add(20)
		for i := 0; i < 20; i++ {
			start := int64(i * 50)
			end := int64((i + 1) * 50)
			go catchByRange(start, end)
		}
		wg.Add(1)
		go removeCatched()

		// Add a timeout to prevent infinite waiting
		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-done:
			// All goroutines completed successfully
			color.Green("All goroutines completed successfully")
			endTime := time.Now()
			color.Cyan("Time taken to remove %d players in Redis: %v\n", totalRemoved, endTime.Sub(startTime))
		case <-time.After(30 * time.Second):
			// Timeout after 30 seconds
			ended_channel <- true
			t.Fatal("Test timed out after 30 seconds")
		}

		close(ch_players)
		close(ended_channel)

		// Verify all players are removed
		remained_players, err := redis_client.GetAllQueueMembersLength(ctx, queue_id)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), remained_players, "All players should be removed from the queue")
		fmt.Printf("Amount of players remained: %d\n", remained_players)
		fmt.Printf("Total players removed: %d\n", totalRemoved)
	})
}
