package main

import (
	"log"
	"net/http"
	"text/template"

	"lyrics/auth"
	dbpkg "lyrics/db"
	handlerspkg "lyrics/handlers"
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
	handlerspkg.SetTemplates(tpl)

	// css chargement

	fs := http.FileServer(http.Dir("../frontend/src/"))
	http.Handle("/css/", http.StripPrefix("/", fs))

	// router http

	http.HandleFunc("/", handlerspkg.AcceuilHandle)
	http.HandleFunc("/forum/index", handlerspkg.ForumIndexHandle)
	http.Handle("/like", auth.RequireAuth(http.HandlerFunc(handlerspkg.LikeHandler)))

	log.Println("🚀 Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// page qui doit etre securiser (login obligatoire)
	// http.Handle("/exemple", auth.RequireAuth(http.HandlerFunc(exempleHandler)))
}
