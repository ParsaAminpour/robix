package main

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/http"
	"os"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/fatih/color"
)

func generateJWTSecretAutomation() {
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

func addUsersToRedisAutomation(amount int) {
	// provide a http call to this route: http://127.0.0.1:8081/matchmake/join
	for i := 0; i < amount; i++ {
		username := petname.Generate(2, "-")
		rating := mathRand.Intn(200) + 10 // Random rating between 0-200
		jsonData, err := json.Marshal(map[string]interface{}{
			"username":            username,
			"match_making_rating": rating,
			"queue_id":            "queue_1",
		})
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			continue
		}
		response, err := http.Post("http://127.0.0.1:8081/matchmake/join", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error adding user %s: %v\n", username, err)
			continue
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			continue
		}
		fmt.Printf("Response body: %s\n", string(body))
		color.Green("Added user %s with rating %d\n", username, rating)
	}
}

func main() {
	generateJWT := flag.NewFlagSet("generate-jwt", flag.ExitOnError)
	addUsers := flag.NewFlagSet("add-users", flag.ExitOnError)

	amount := addUsers.Int("amount", 1, "Number of users to add")

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
		addUsersToRedisAutomation(*amount)
	default:
		fmt.Println("Expected 'generate-jwt' or 'add-users' subcommands")
		os.Exit(1)
	}
}
