package main

import (
	"log"

	kvs "github.com/defoeam/herd/internal"
)

func main() {
	if err := kvs.StartGRPCServer(true); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
