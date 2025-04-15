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

	api "github.com/ParsaAminpour/robix/backend/matchmaking/handler"
	"github.com/ParsaAminpour/robix/backend/matchmaking/models"
	"github.com/ParsaAminpour/robix/backend/matchmaking/redis"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	redis_client  *redis.RedisClient
	ctx           context.Context
	wg            sync.WaitGroup
	once          sync.Once
	queue_id      string = "queue_1"
	upgrader      *websocket.Upgrader
	Clients       map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	mu            sync.Mutex
	joinedMessage string = "You Joined the Game %s"
)

func setup() {
	once.Do(func() {
		var err error
		redis_client, err = redis.ConnectToRedis(ctx)
		if err != nil {
			log.Fatalf("failed to connect to redis: %v", err)
		}
		if redis_client == nil {
			log.Fatal("Redis client is nil after initialization")
		}

		upgrader = &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		fmt.Println("Redis client and ws upgrader initialized")
	})
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
	fmt.Printf("You are in Handle Websocket with ID: %s\n", username)
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

	setup()
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
	go removeByRange(0, 100)
	go removeByRange(100, 200)
	go removeByRange(200, 300)

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
		matchmaking.GET("/queue", endpointHandler(api.GetMembersFromQueue))
		matchmaking.GET("/clients", func(c echo.Context) error {
			jsonData, err := json.Marshal(Clients)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusOK, map[string]interface{}{"message": jsonData})
		})
		matchmaking.POST("/join", joinQueueEndpointHandler(api.AddToQueue))
		matchmaking.POST("/leave", endpointHandler(api.RemoveFromQueue))
	}

	if err := e.Start(":8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}
}
