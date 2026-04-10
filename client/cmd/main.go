package main

import "client"

func main() {
	cl := client.NewClient("anton")
	cl.Start()
}
