package handlers

import (
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template

func SetTemplates(t *template.Template) {
	tpl = t
}

func AcceuilHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "accueil.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page d'accueil", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

func ForumIndexHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page du forum", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}

func ProfileHandle(w http.ResponseWriter, r *http.Request) {
	if tpl == nil {
		http.Error(w, "templates non initialisés", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "profile.html", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page de profil", http.StatusInternalServerError)
		log.Printf("Erreur template: %v", err)
	}
}
