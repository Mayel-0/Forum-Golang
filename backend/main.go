package main

import (
	"log"
	"text/template"

	dbpkg "lyrics/db"
	tplpkg "lyrics/template"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var err error
var tpl *template.Template
var db *gorm.DB

func main() {
	if err = godotenv.Load("env/.env"); err == nil {
		log.Println("✅ Variables d'environnement chargées depuis env/.env")
	} else {
		log.Fatalf("Erreur chargement .env: %v", err)
	}

	if db, err = dbpkg.ConnectDB(); err != nil {
		log.Fatalf("Erreur de connexion à la DB: %v", err)
	}

	tpl, err = tplpkg.ParseTemplates()
	if err != nil {
		log.Fatal("erreur template", err)
	}
}
