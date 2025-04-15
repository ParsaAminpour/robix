package internal

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	mu            sync.Mutex
	Clients       map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	joinedMessage string                     = "You Joined the Game %s"
	upgrader      *websocket.Upgrader
)

func HandleWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer ws.Close()

	username := c.QueryParam("username")
	fmt.Printf("You are in Handle Websocker with ID: %s\n", username)
	mu.Lock()
	Clients[username] = ws
	mu.Unlock()
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			// the connection lost
			mu.Lock()
			delete(Clients, username)
			mu.Unlock()
			fmt.Printf("User %s has been disconnected", username)
			break
		}
	}
	return nil
}
