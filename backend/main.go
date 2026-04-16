package main

import (
	"log"
	"net/http"
	"text/template"

	dbpkg "lyrics/db"
	tplpkg "lyrics/template"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var err error
var tpl *template.Template
var db *gorm.DB

func acceuilHandle(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "accueil.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page d'accueil", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

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

	// css chargement

	fs := http.FileServer(http.Dir("../frontend/src/"))
	http.Handle("/css/", http.StripPrefix("/", fs))

	// router http

	http.HandleFunc("/", acceuilHandle)

	log.Println("🚀 Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
