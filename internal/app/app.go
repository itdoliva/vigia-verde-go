package app

import (
	"context"
	"log"
	"net/http"

	"vigia-verde-go/internal/handlers"
	"vigia-verde-go/internal/platform"
	"vigia-verde-go/internal/repository"
	service "vigia-verde-go/internal/service"
)

func Run() error {
	ctx := context.Background()

	fsClient, cleanup, err := platform.SetupDatabaseConnection(ctx)
	if err != nil {
		return err
	}
	defer cleanup()

	treeRepo := repository.NewTreeEventRepository(fsClient)
	treeService := service.NewTreeEventService(treeRepo)
	treeHandler := handlers.NewTreeEventHandler(treeService)

	mux := http.NewServeMux()
	treeHandler.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("Servidor HTTP ouvindo em %s ...", addr)

	return http.ListenAndServe(addr, mux)
}
