package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"

	"vigia-verde-go/internal/core/config"
)

func main() {
	// 1) Carrega .env (ajuste o caminho dependendo de onde você roda o go run)
	if err := config.LoadEnv("../../.env"); err != nil {
		log.Fatalf("erro ao carregar .env: %v", err)
	}

	// 2) Cria contexto base
	ctx := context.Background()

	// 3) Inicializa Firebase App
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("erro ao inicializar Firebase App: %v", err)
	}

	// 4) Cria Firestore client
	fsClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("erro ao criar Firestore client: %v", err)
	}
	defer fsClient.Close()

	log.Println("Conectado ao Firestore com sucesso! ✅")

}
