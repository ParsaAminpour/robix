package tests

import (
	"context"
	"testing"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/internal"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()

func TestRedisConnectionAndPing(t *testing.T) {
	redis_client, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	ping, ping_err := redis_client.Ping(context.Background()).Result()
	assert.NoError(t, ping_err)
	assert.Equal(t, ping, "PONG")
}

func TestRedisSetAndGetKeyValue(t *testing.T) {
	_, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	err = redis.SetValue(ctx, "name", "john")
	assert.NoError(t, err)

	fetched_value, fetched_err := redis.GetValue(ctx, "name")
	assert.NoError(t, fetched_err)
	assert.Equal(t, fetched_value, "john")
}

func TestRedisZAdd(t *testing.T) {
	_, err := redis.ConnectToRedis(ctx)
	assert.NoError(t, err)

	player := models.Player{}
	another_player := models.Player{}
	player = player.NewPlayer("john", uuid.NewString(), 0)
	another_player = another_player.NewPlayer("Elliot", uuid.NewString(), 0)

	t.Run("Test ZAdd", func(t *testing.T) {
		err = redis.AddOrUpdatePlayerQueueMMR(ctx, &player)
		assert.NoError(t, err)
	})

	t.Run("Test ZUpdate", func(t *testing.T) {
		err = redis.AddOrUpdatePlayerQueueMMR(ctx, &player)
		assert.NoError(t, err)
	})

	t.Run("Test ZRem", func(t *testing.T) {
		err = redis.AddOrUpdatePlayerQueueMMR(ctx, &another_player)
		err = redis.RemovePlayerFromQueueMMR(ctx, &player)
		assert.NoError(t, err)

		players, err := redis.GetAllQueueMemeber(ctx, another_player.QueueID)
		assert.NoError(t, err)
		assert.Equal(t, players[0].ID, another_player.ID)
		assert.Equal(t, players[0].Score, another_player.MatchmakingRating)
		assert.Len(t, players, 1) // only one player is in queue
	})
}
