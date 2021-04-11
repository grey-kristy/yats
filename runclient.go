package main

import (
	"flag"
	"yats/client"
)

func main() {
	port := "localhost:9000"
	n := flag.Int("number", 1, "number of data packets")
	flag.Parse()
	client.Client(port, *n)
}
