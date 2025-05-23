package handler

import (
	"context"
	"net/http"
	"sync"

	models "github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	mu = &sync.Mutex{}
)

func AddToQueue(c echo.Context, ctx context.Context, redis_client *redis.RedisClient, clients map[string]*websocket.Conn) error {
	mu.Lock()
	defer mu.Unlock()

	player := new(models.Player)
	if err := c.Bind(player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if clients[player.Username] == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User is not connected to the websocket"})
	}
	player = player.NewPlayer(player.Username, player.QueueID, player.MatchmakingRating)

	if err := redis_client.AddOrUpdatePlayerQueueMMR(ctx, player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Player added to queue"})
}

func RemoveFromQueue(c echo.Context, ctx context.Context, redis_client *redis.RedisClient) error {
	mu.Lock()
	defer mu.Unlock()

	player := new(models.Player)
	if err := c.Bind(player); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := redis_client.RemovePlayerFromQueueMMR(ctx, player); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Player removed from queue"})
}

func GetMembersFromQueue(c echo.Context, ctx context.Context, redis_client *redis.RedisClient) error {
	mu.Lock()
	defer mu.Unlock()

	// queue_id := c.Param("queue_id")
	members, err := redis_client.GetAllQueueMembers(ctx, "queue_1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, members)
}
