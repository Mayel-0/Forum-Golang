package main

import (
	"log"

	db "lyrics/db"

	"github.com/joho/godotenv"
)

var err error

func main() {
	if err = godotenv.Load("env/.env"); err == nil {
		log.Println("✅ Variables d'environnement chargées depuis env/.env")
	} else {
		log.Fatalf("Erreur chargement .env: %v", err)
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Erreur de connexion à la DB: %v", err)
	}

}
