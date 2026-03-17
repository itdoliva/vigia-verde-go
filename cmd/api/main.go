package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"vigia-verde-go/internal/infrastructure/firebase"
	handlerEvent "vigia-verde-go/internal/interface/http/event"
	config "vigia-verde-go/internal/shared"
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

	eventHandler := handlerEvent.SetupModule(fsClient)
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
