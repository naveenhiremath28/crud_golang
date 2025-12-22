package main

import (
	"log"
	"practise/go_fiber/internal/containers"
)

func main() {
	container, err := containers.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	
	// Start server (blocking)
	if err := container.Invoke(containers.StartServer); err != nil {
		log.Fatal(err)
	}
}
