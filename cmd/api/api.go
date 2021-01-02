package main

import "api/internal/app/driver/webserver"

func main() {
	server := webserver.NewServer()
	server.Run()
}