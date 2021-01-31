package main

import (
	"api/internal/app"
	"api/internal/app/driver/grpcserver"
	"fmt"
	"os"
)

func main() {
	if err := setup(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server, err := grpcserver.NewServer()
	if err != nil {
		fmt.Printf("server initialize error: %s", err)
		os.Exit(1)
	}

	server.Run()
}

func setup() error {
	if err := app.InitDI(); err != nil {
		return fmt.Errorf("initialize DI error: %w", err)
	}
	if err := app.InitConfig(); err != nil {
		return fmt.Errorf("initialize Config error: %w", err)
	}
	if err := app.InitLogger(); err != nil {
		return fmt.Errorf("initialize Logger error: %w", err)
	}
	return nil
}
