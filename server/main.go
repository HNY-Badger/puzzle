package main

import (
	"bufio"
	"context"
	"fmt"
	"godb/application"
	"os"
	"os/signal"
	"strings"
)

func main() {
	fmt.Println("Initializing server...")
	app := application.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fmt.Print("Enter port for the app (3000) : ")
	reader := bufio.NewReader(os.Stdin)
	var port string
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, string([]byte{13}), "", -1)
	if len(text) == 0 {
		port = "3000"
	} else {
		port = text
	}

	err := app.Start(ctx, port)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
