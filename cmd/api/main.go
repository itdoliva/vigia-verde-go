package main

import (
	"context"
	"log"
	"net/http"
	"vigia-verde-go/internal/core/config"
	handlers "vigia-verde-go/internal/handlers"
	platform "vigia-verde-go/internal/platform"
	"vigia-verde-go/internal/repository"
	service "vigia-verde-go/internal/service"
)

func main() {
	config.LoadEnv("../../.env")

	ctx := context.Background()
	fsClient, cleanup, err := platform.SetupDatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	treeRepo := repository.NewTreeEventRepository(fsClient)
	treeService := service.NewTreeEventService(treeRepo)
	treeHandler := handlers.NewTreeEventHandler(treeService)

	mux := http.NewServeMux()
	treeHandler.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("Servidor HTTP ouvindo em %s ...", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("erro ao subir servidor HTTP: %v", err)
	}
	defer cleanup()
}
