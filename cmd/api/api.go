package main

import (
	"api/internal/app/driver/grpcserver"
	"api/internal/app/util/di"
	"api/internal/app/util/di/container/prod"
)

func main() {
	// TODO: env
	di.SetDi(&prod.Container{})

	server := grpcserver.NewServer()
	server.Setup()
	server.Run()
}
