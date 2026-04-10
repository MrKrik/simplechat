package main

import server "tcp-server"

func main() {
	server := server.NewServer("localhost:3000")
	server.Start()
}
