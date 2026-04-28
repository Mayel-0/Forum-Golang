package main

import (
	"log"
	"net/http"
	"text/template"

	"lyrics/auth"
	dbpkg "lyrics/db"
	handlerspkg "lyrics/handlers"
	tplpkg "lyrics/template"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var err error
var tpl *template.Template
var db *gorm.DB
var store = sessions.NewCookieStore([]byte("votre-cle-secrete-tres-longue-et-aleatoire-32-octets-minimum"))

func main() {
	auth.SetStore(store)
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
	assetsFS := http.FileServer(http.Dir("../frontend/assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetsFS))

	// router http

	http.HandleFunc("/", handlerspkg.AcceuilHandle)

	http.HandleFunc("/forum/index", handlerspkg.ForumIndexHandle)
	http.HandleFunc("/auth/login", handlerspkg.LoginHandle)
	http.HandleFunc("/auth/register", handlerspkg.RegisterHandle)
	http.HandleFunc("/auth/logout", handlerspkg.LogoutHandle)

	http.Handle("/profile", auth.RequireAuth(http.HandlerFunc(handlerspkg.ProfileHandle)))
	http.Handle("/profile/modify", auth.RequireAuth(http.HandlerFunc(handlerspkg.ProfileHandle)))
	http.Handle("/like/add", auth.RequireAuth(http.HandlerFunc(handlerspkg.LikeHandlerAdd)))
	http.Handle("/like/rm", auth.RequireAuth(http.HandlerFunc(handlerspkg.LikeHandlerRm)))

	http.Handle("/poste/create", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteCreateHandler)))
	http.Handle("/poste/modifier", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteModifierHandle)))
	http.Handle("/poste/supprimer", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteDeleteHandler)))

	log.Println("🚀 Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// page qui doit etre securiser (login obligatoire)
	// http.Handle("/exemple", auth.RequireAuth(http.HandlerFunc(exempleHandler)))
}
