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
	"time"

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
)

func endpointHandler(_handler func(c echo.Context, ctx context.Context, redis_client *redis.RedisClient) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return _handler(c, ctx, redis_client)
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx = context.Background()

	// Setup Signal Handling
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize Redis
	var err error
	redis_client, err = redis.ConnectToRedis(ctx)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	if redis_client == nil {
		log.Fatal("Redis client is nil after initialization")
	}

	ch_players := make(chan []models.AbstractPlayer, 10) // todo: use dynamic buffered channel

	// writer goroutine
	wg.Add(1)
	go func() {
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
				fmt.Println("getting players from redis")
				players, err := redis_client.GetPlayersByRatingRange(ctx, "queue_1", 0, 100, 5)
				if err != nil {
					log.Printf("Error getting players: %v", err)
					continue
				}
				if len(players) > 0 {
					color.Green("found players in writer goroutine: %v", players)
					ch_players <- players
				}
			}
			time.Sleep(time.Second * 1)
		}
		color.Red("Writer goroutine finished")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-signalChan:
			color.Red("Received interrupt signal from reader goroutine, shutting down reader goroutine...")
			close(ch_players)
			os.Exit(0)
			break
		case players := <-ch_players:
			color.Green("catch players in reader goroutine: %v", players)
			if err := redis_client.RemoveBatchPlayersFromQueueMMR(ctx, players); err != nil {
				log.Printf("Error removing players: %v", err)
			}
			color.Green("removed players in reader goroutine: %v", players)
		}
	}()

	go func() {
		wg.Wait()
		color.Red("All goroutines finished")
		os.Exit(0)
	}()

	e := echo.New()

	// Global middleware
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
	// matchmaking.Use(customMiddleware.AuthMiddleware(authConfig))
	{
		matchmaking.GET("/queue", endpointHandler(api.GetMembersFromQueue))
		matchmaking.POST("/join", endpointHandler(api.AddToQueue))
		matchmaking.POST("/leave", endpointHandler(api.RemoveFromQueue))
	}

	if err := e.Start(":8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}
}
