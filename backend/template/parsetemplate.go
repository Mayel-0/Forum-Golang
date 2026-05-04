package template

import "text/template"

var tpl *template.Template
var err error

func ParseTemplates() (*template.Template, error) {
	tpl, err = template.ParseFiles(
		"../frontend/src/html/accueil.html",
		"../frontend/src/html/profile/profile.html",
		"../frontend/src/html/forum/index.html",
		"../frontend/src/html/partials/navbar.html",
		"../frontend/src/html/partials/footer.html",
		"../frontend/src/html/auth/login.html",
		"../frontend/src/html/auth/register.html",
		"../frontend/src/html/forum/post/post-create.html",
		"../frontend/src/html/forum/post/post-edit.html",
	)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
