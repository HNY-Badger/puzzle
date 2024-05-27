package main

import (
	"context"
	"fmt"
	"godb/application"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("Initializing server...")
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
