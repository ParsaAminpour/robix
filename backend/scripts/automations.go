package main

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	mathRand "math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var wg sync.WaitGroup
var signalChan chan os.Signal

func generateJWTSecretAutomation() {
	fmt.Println("hahha")
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := cryptoRand.Read(bytes); err != nil {
		fmt.Printf("Error generating random bytes: %v\n", err)
		return
	}

	// Encode as base64
	secret := base64.StdEncoding.EncodeToString(bytes)
	color.Green("Generated JWT Secret: %s\n", secret)
}

func addUserToRedisAutomation() {
	// provide a http call to this route: http://127.0.0.1:8081/matchmake/join
	username := petname.Generate(2, "-")
	rating := mathRand.Intn(200) + 10 // Random rating between 0-200
	jsonData, err := json.Marshal(map[string]interface{}{
		"username":            username,
		"match_making_rating": rating,
		"queue_id":            "queue_1",
	})
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
	}

	route := fmt.Sprintf("ws://127.0.0.1:8081/matchmake/ws?username=%s", username)
	conn, _, err := websocket.DefaultDialer.Dial(route, nil)
	if err != nil {
		log.Fatal("error while trying to stablish a ws connection")
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("received message: %s", message)
		}
	}()

	response, err := http.Post("http://127.0.0.1:8081/matchmake/join", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error adding user %s: %v\n", username, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
	}
	fmt.Printf("Response body: %s\n", string(body))
	color.Green("Added user %s with rating %d\n", username, rating)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	go func() {
		for signal := range signalChan {
			color.Red("Received %v signal from writer goroutine, shutting down writer goroutine...", signal)
			os.Exit(0)
			wg.Done()
		}
	}()

	for {
		tickerTime := <-ticker.C
		sendMessage := fmt.Sprintf("hi - %s", tickerTime.String())
		err = conn.WriteMessage(websocket.TextMessage, []byte(sendMessage))
		if err != nil {
			log.Println("write error:", err)
		} else {
			log.Println("send message:", sendMessage)
		}
	}
}

func connectToWebsocket(amount int) {
	for i := 0; i < amount; i++ {
		username := petname.Generate(1, "")
		route := fmt.Sprintf("ws://127.0.0.1:8081/matchmake/ws?username=%s", username)

		conn, _, err := websocket.DefaultDialer.Dial(route, nil)
		if err != nil {
			log.Fatal("error while trying to stablish a ws connection")
		}
		defer conn.Close()

		fmt.Printf("Response body: %s\n", conn.RemoteAddr().String())
		color.Blue("Added user %s to websocket\n", username)
	}
}

func main() {
	generateJWT := flag.NewFlagSet("generate-jwt", flag.ExitOnError)
	addUsers := flag.NewFlagSet("add-users", flag.ExitOnError)
	wsConnection := flag.NewFlagSet("add-to-ws", flag.ExitOnError)

	addUseramount := addUsers.Int("amount", 1, "Number of users to add")
	wsConnectionAmount := wsConnection.Int("amount", 1, "Number of users to add ws")

	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	if len(os.Args) < 2 {
		fmt.Println("Expected 'generate-jwt' or 'add-users' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate-jwt":
		generateJWT.Parse(os.Args[2:])
		generateJWTSecretAutomation()

	case "add-users":
		addUsers.Parse(os.Args[2:])
		wg.Add(*addUseramount)
		defer wg.Done()

		for i := 0; i < *addUseramount; i++ {
			go addUserToRedisAutomation()
			time.Sleep(time.Millisecond * 10)
			fmt.Println("----------------")
		}

	case "add-to-ws":
		wsConnection.Parse(os.Args[2:])
		connectToWebsocket(*wsConnectionAmount)

	default:
		fmt.Println("Expected 'generate-jwt' or 'add-users' subcommands")
		os.Exit(1)
	}

	go func() {
		wg.Wait()
		os.Exit(1)
	}()
}
