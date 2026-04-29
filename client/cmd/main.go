package main

import (
	client "client/iternal"
)

func main() {
	cl := client.NewClient("anton")
	cl.Start()
}
