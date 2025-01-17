package main

import (
	"flag"
	"log"

	kvs "github.com/defoeam/herd/internal"
)

func main() {
	useLogging := flag.Bool("useLogging", false, "Enable logging")
	useSecurity := flag.Bool("useSecurity", false, "Enable security")

	flag.Parse()

	if err := kvs.StartGRPCServer(*useLogging, *useSecurity); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
