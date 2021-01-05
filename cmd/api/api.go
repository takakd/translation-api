package main

import "api/internal/app/driver/grpcserver"

func main() {
	server := grpcserver.NewServer()
	server.Setup()
	server.Run()
}
