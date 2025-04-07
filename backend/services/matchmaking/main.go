package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	api "github.com/ParsaAminpour/robix/backend/matchmaking/handler"
	GenerateJWTSecret "github.com/ParsaAminpour/robix/backend/matchmaking/middleware"
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
	mu           sync.Mutex
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

	// Initialize Redis
	var err error
	redis_client, err = redis.ConnectToRedis(ctx)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	if redis_client == nil {
		log.Fatal("Redis client is nil after initialization")
	}

	// remove this line after generating the JWT secret
	GenerateJWTSecret.GenerateJWTSecret()

	// Initialize matchmaking workers
	// wg.Add(5)
	// channel := make(chan []models.AbstractPlayer)
	// for i := 0; i < 5; i++ {
	// 	go internal.MatchMakeByRatingRange(&wg, channel, ctx, redis_client)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(channel)
	// }()

	e := echo.New()

	// Configure auth middleware
	// authConfig := customMiddleware.AuthConfig{
	// 	JWTSecret:      os.Getenv("JWT_SECRET"),
	// 	UserServiceURL: os.Getenv("USER_SERVICE_URL"),
	// }

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
