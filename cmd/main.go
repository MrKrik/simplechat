package main

import "tcp-server/server"

func main() {
	server := server.NewServer(":3000")
	server.Start()
}
