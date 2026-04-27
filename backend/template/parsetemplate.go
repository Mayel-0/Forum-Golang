package template

import "text/template"

var tpl *template.Template
var err error

func ParseTemplates() (*template.Template, error) {
	tpl, err = template.ParseFiles(
		"../frontend/src/html/accueil.html",
		"../frontend/src/html/forum/index.html",
		"../frontend/src/html/partials/navbar.html",
		"../frontend/src/html/partials/footer.html",
		"../frontend/src/html/auth/login.html",
		"../frontend/src/html/auth/register.html",
	)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
