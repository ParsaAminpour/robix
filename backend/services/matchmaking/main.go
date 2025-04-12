package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	api "github.com/ParsaAminpour/robix/backend/matchmaking/handler"
	"github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	redis_client *redis.RedisClient
	ctx          context.Context
	wg           sync.WaitGroup
	queue_id     string = "queue_1"
)

func endpointHandler(_handler func(c echo.Context, ctx context.Context, redis_client *redis.RedisClient) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return _handler(c, ctx, redis_client)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx = context.Background()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var err error
	redis_client, err = redis.ConnectToRedis(ctx)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	if redis_client == nil {
		log.Fatal("Redis client is nil after initialization")
	}

	// TODO: use dynamic buffered channel
	ch_players := make(chan []models.AbstractPlayer, 100)

	removeByRange := func(start, end int64) {
		defer wg.Done()
	WRITER_LOOP:
		for {
			select {
			case <-signalChan:
				color.Red("Received interrupt signal from writer goroutine, shutting down writer goroutine...")
				close(ch_players)
				os.Exit(0)
				break WRITER_LOOP
			default:
				players, err := redis_client.GetPlayersByRatingRange(ctx, queue_id, start, end, 5)
				if err != nil {
					log.Printf("Error getting players: %v", err)
					continue
				}
				if len(players) > 0 {
					fmt.Printf("found players in writer goroutine range %d-%d: %d\n", start, end, len(players))
					ch_players <- players
				}
			}
		}
		color.Red("Writer goroutine finished")
	}

	readFromChannel := func() {
		defer wg.Done()
	READER_LOOP:
		for {
			select {
			case <-signalChan:
				color.Red("Received interrupt signal from reader goroutine, shutting down reader goroutine...")
				close(ch_players)
				os.Exit(0)
				break READER_LOOP
			case players := <-ch_players:
				if err := redis_client.RemoveBatchPlayersFromQueueMMR(ctx, players); err != nil {
					log.Printf("Error removing players: %v", err)
				}
				color.Red("removed players in reader goroutine:")
			}
		}
	}

	// writer goroutine
	wg.Add(3)
	go removeByRange(0, 100)
	go removeByRange(100, 200)
	go removeByRange(200, 300)

	// reader goroutine
	wg.Add(3)
	go readFromChannel()
	go readFromChannel()
	go readFromChannel()

	go func() {
		wg.Wait()
		color.Red("All goroutines finished")
		os.Exit(0)
	}()

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/metrics"
		},
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value, _ := c.Get("customValueFromContext").(int)
			color.Green("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	matchmaking := e.Group("/matchmake")
	{
		matchmaking.GET("/queue", endpointHandler(api.GetMembersFromQueue))
		matchmaking.POST("/join", endpointHandler(api.AddToQueue))
		matchmaking.POST("/leave", endpointHandler(api.RemoveFromQueue))
	}

	if err := e.Start(":8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}
}
