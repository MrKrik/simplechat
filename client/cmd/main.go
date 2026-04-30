package main

import (
	client "client/iternal"
)

func main() {
	cl := client.NewClient()
	cl.Start()
}
