package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"vigia-verde-go/internal/event"
	"vigia-verde-go/internal/platform/config"
	"vigia-verde-go/internal/platform/firebase"
)

func main() {
	config.LoadEnv(".env")

	ctx := context.Background()
	fsClient, err := firebase.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	mux := http.NewServeMux()

	eventHandler := event.SetupModule(fsClient)
	eventHandler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Servidor HTTP ouvindo em %s ...", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("erro ao subir servidor HTTP: %v", err)
	}
}
