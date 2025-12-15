package main

import (
	"context"
	"log"
	"vigia-verde-go/internal/core/config"
	platform "vigia-verde-go/internal/platform"
)

func main() {
	config.LoadEnv("../../.env")

	ctx := context.Background()
	_, cleanup, err := platform.SetupDatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	defer cleanup()
}
