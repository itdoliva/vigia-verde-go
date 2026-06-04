package main

import (
	"context"
	"log"
	"net/http"
	"os"
	_ "vigia-verde-go/cmd/api/docs"
	"vigia-verde-go/internal/infrastructure/firebase"
	eventModule "vigia-verde-go/internal/interface/http/event"
	config "vigia-verde-go/internal/shared"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Vigia Verde API Docs
// @version 0.1.0
// @contact.name Arthur de Oliveira & Italo de Oliveira
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host vigia-verde-514678222430.southamerica-east1.run.app
// @BasePath /
func main() {
	config.LoadEnv(".env")

	ctx := context.Background()
	fsClient, err := firebase.NewFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	mux := http.NewServeMux()

	mux.Handle("/docs/", httpSwagger.WrapHandler)

	eventHandler := eventModule.SetupModule(fsClient)
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
