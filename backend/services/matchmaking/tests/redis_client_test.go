package tests

import (
	"context"
	"testing"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var redis_client *redis.RedisClient
var queue_id string = uuid.NewString()

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

	player := models.Player{}
	another_player := models.Player{}
	player = player.NewPlayer("john", queue_id, 0)
	another_player = another_player.NewPlayer("Elliot", queue_id, 0)

	t.Run("Test ZAdd", func(t *testing.T) {
		err = redis_client.AddOrUpdatePlayerQueueMMR(ctx, &player)
		assert.NoError(t, err)
	})

	t.Run("Test ZUpdate", func(t *testing.T) {
		err = redis_client.AddOrUpdatePlayerQueueMMR(ctx, &player)
		assert.NoError(t, err)
	})

	t.Run("Test ZRem", func(t *testing.T) {
		err = redis_client.AddOrUpdatePlayerQueueMMR(ctx, &another_player)
		err = redis_client.RemovePlayerFromQueueMMR(ctx, &player)
		assert.NoError(t, err)

		players, err := redis_client.GetAllQueueMembers(ctx, player.QueueID)
		assert.NoError(t, err)
		assert.Equal(t, players[0].ID, another_player.ID)
		assert.Equal(t, players[0].Score, another_player.MatchmakingRating)
		assert.Len(t, players, 1) // only one player is in queue
	})
}
