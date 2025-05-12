package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	matchmaking_api "github.com/ParsaAminpour/robix/backend/matchmaking/handler"
	"github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	consul_api "github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	redis_client *redis.RedisClient
	ctx          context.Context
	wg           sync.WaitGroup
	once         sync.Once
	upgrader     *websocket.Upgrader
	Clients      map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	mu           sync.Mutex
	consulClient *consul_api.Client
)

// TODO: Avoid hardcoding Consul Service Registry configs
const (
	joinedMessage string = "You Joined the Game %s"
	serviceId     string = "matchmaking_service"
	queue_id      string = "queue_1"
	port          int    = 8081
	address       string = "127.0.0.1"
)

func setup() error {
	once.Do(func() {
		var err error
		redis_client, err = redis.ConnectToRedis(ctx)
		if err != nil {
			log.Fatalf("failed to connect to redis: %v", err)
		}
		if redis_client == nil {
			log.Fatal("Redis client is nil after initialization")
		}

		// Initialize Consul client
		// consulConfig := consul_api.DefaultConfig()
		// consulClient, err = consul_api.NewClient(consulConfig)
		// if err != nil {
		// 	log.Fatal("error during initializing Consul client:", err)
		// }

		// // Register service with Consul
		// registration := &consul_api.AgentServiceRegistration{
		// 	ID:      serviceId,
		// 	Name:    "matchmaking-service",
		// 	Port:    port,
		// 	Address: address,
		// 	Check: &consul_api.AgentServiceCheck{
		// 		HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
		// 		Interval: "10s",
		// 		Timeout:  "5s",
		// 	},
		// }

		// if err := consulClient.Agent().ServiceRegister(registration); err != nil {
		// 	log.Fatal("error registering service with Consul:", err)
		// }

		upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		fmt.Println("Redis client, Consul client, and ws upgrader initialized")
	})
	return nil
}

func endpointHandler(_handler func(c echo.Context, ctx context.Context, redis_client *redis.RedisClient) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return _handler(c, ctx, redis_client)
	}
}

func joinQueueEndpointHandler(_handler func(c echo.Context, ctx context.Context, redis_client *redis.RedisClient, _ map[string]*websocket.Conn) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return _handler(c, ctx, redis_client, Clients)
	}
}

func HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer ws.Close()

	username := c.QueryParam("username")
	mu.Lock()
	Clients[username] = ws
	mu.Unlock()

	if err := ws.WriteMessage(websocket.TextMessage, []byte("Hi, this is test ws")); err != nil {
		fmt.Println("Error Occurred in sending ws message")
	}
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// the connection lost
			mu.Lock()
			delete(Clients, username)
			mu.Unlock()
			redis_client.RemovePlayerByUsername(ctx, username)
			color.Red("User %s has been disconnected\n", username)
			break
		}
	}
	return nil
}

func sendNotificationToPlayers(usernames []string, message string) error {
	for _, username := range usernames {
		if ws_conn, exist := Clients[username]; exist {
			if err := ws_conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				return err
			}
			// @audit-info should this deletion applied also if there was an error?
			delete(Clients, username)
		}
	}
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx = context.Background()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if err := setup(); err != nil {
		log.Fatal("Failed to setup services:", err)
	}
	// @audit-info use dynamic buffered channel
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
					// @audit what if the buffer is full?
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
				// Send Notification to players
				var players_id []string
				for _, player := range players {
					players_id = append(players_id, player.ID)
				}
				sendNotificationToPlayers(players_id, joinedMessage)

				color.Red("removed players in reader goroutine")
			}
		}
	}
	// writer goroutine
	wg.Add(3)
	ranges := []struct{ start, end int64 }{
		{0, 100},
		{100, 200},
		{200, 300},
	}
	for _, r := range ranges {
		go removeByRange(r.start, r.end)
	}

	// reader goroutine
	numReaders := 6
	wg.Add(numReaders)
	for i := 0; i < numReaders; i++ {
		go readFromChannel()
	}

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
		matchmaking.GET("/ws", HandleWebSocket)
		matchmaking.GET("/queue", endpointHandler(matchmaking_api.GetMembersFromQueue))
		matchmaking.GET("/clients", func(c echo.Context) error {
			jsonData, err := json.Marshal(Clients)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusOK, map[string]interface{}{"message": jsonData})
		})
		matchmaking.POST("/join", joinQueueEndpointHandler(matchmaking_api.AddToQueue))
		matchmaking.POST("/leave", endpointHandler(matchmaking_api.RemoveFromQueue))
	}

	// Add health check endpoint for Consul
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Start server in a goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	<-signalChan
	log.Println("Received shutdown signal")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Deregister service from Consul
	// if err := consulClient.Agent().ServiceDeregister(serviceId); err != nil {
	// 	log.Printf("Error deregistering service from Consul: %v", err)
	// }

	// Shutdown Echo server
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Server gracefully stopped")
}
