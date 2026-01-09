package main

import (
	"log"
	"vigia-verde-go/internal/app"
	"vigia-verde-go/internal/core/config"
)

func main() {
	config.LoadEnv("../../.env")

	if err := app.Run(); err != nil {
		log.Fatalf("erro ao subir servidor HTTP: %v", err)
	}
}
