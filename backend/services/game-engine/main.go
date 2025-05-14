package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
)

var (
	ctx = context.Background()
)

func main() {
	e := echo.New()

	_, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println("Game engine started")

	e.Start(":8080")
}
