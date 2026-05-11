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
var tpl map[string]*template.Template
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

	mux := http.NewServeMux()

	// CSS & Assets
	fs := http.FileServer(http.Dir("../frontend/src/"))
	mux.Handle("/css/", http.StripPrefix("/", fs))
	assetsFS := http.FileServer(http.Dir("../frontend/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetsFS))

	// Routes publiques
	mux.HandleFunc("/", handlerspkg.ForumIndexHandle)
	mux.HandleFunc("GET /category/{slug}", handlerspkg.CategoryHandle)
	mux.HandleFunc("GET /category/{slug}/{subcategory}", handlerspkg.SubCategoryHandle)
	mux.HandleFunc("GET /auth/login", handlerspkg.LoginHandle)
	mux.HandleFunc("POST /auth/login", handlerspkg.LoginHandle)
	mux.HandleFunc("GET /auth/register", handlerspkg.RegisterHandle)
	mux.HandleFunc("POST /auth/register", handlerspkg.RegisterHandle)
	mux.HandleFunc("GET /auth/logout", handlerspkg.LogoutHandle)

	// Route dynamique par slug
	mux.HandleFunc("GET /p/{slug...}", handlerspkg.PostShowHandle)
	// Routes protégées
	mux.Handle("GET /profile", auth.RequireAuth(http.HandlerFunc(handlerspkg.ProfileHandle)))
	mux.Handle("POST /profile/modify", auth.RequireAuth(http.HandlerFunc(handlerspkg.ProfileHandle)))
	mux.Handle("POST /profile/avatar", auth.RequireAuth(http.HandlerFunc(handlerspkg.UploadAvatarHandler)))
	mux.Handle("POST /profile/banniere", auth.RequireAuth(http.HandlerFunc(handlerspkg.UploadBanniereHandler)))

	mux.Handle("POST /like/add", auth.RequireAuth(http.HandlerFunc(handlerspkg.LikeHandlerAdd)))
	mux.Handle("POST /like/rm", auth.RequireAuth(http.HandlerFunc(handlerspkg.LikeHandlerRm)))

	mux.Handle("GET /post/create", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteCreateHandler)))
	mux.Handle("POST /post/create", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteCreateHandler)))
	mux.Handle("POST /post/modifier", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteModifierHandle)))
	mux.Handle("POST /post/supprimer", auth.RequireAuth(http.HandlerFunc(handlerspkg.PosteDeleteHandler)))

	mux.Handle("POST /comment/create", auth.RequireAuth(http.HandlerFunc(handlerspkg.AddCommentHandler)))
	mux.Handle("POST /comment/delete", auth.RequireAuth(http.HandlerFunc(handlerspkg.DeleteCommentHandler)))

	log.Println("🚀 Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
