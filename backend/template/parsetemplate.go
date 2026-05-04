package template

import "text/template"

const base = "../frontend/src/html/"

var partials = []string{
	base + "partials/template.html",
	base + "partials/navbar.html",
	base + "partials/footer.html",
}

func parsePage(pagePath string) (*template.Template, error) {
	files := append([]string{pagePath}, partials...)
	return template.ParseFiles(files...)
}

func ParseTemplates() (map[string]*template.Template, error) {
	pages := map[string]string{
		"accueil.html":     base + "accueil.html",
		"index.html":       base + "forum/index.html",
		"post.html":        base + "forum/post/post.html",
		"post-create.html": base + "forum/post/post-create.html",
		"post-edit.html":   base + "forum/post/post-edit.html",
		"login.html":       base + "auth/login.html",
		"register.html":    base + "auth/register.html",
		"profile.html":     base + "profile/profile.html",
	}

	templates := make(map[string]*template.Template)
	for name, path := range pages {
		t, err := parsePage(path)
		if err != nil {
			return nil, err
		}
		templates[name] = t
	}
	return templates, nil
}
