package template

import (
	"lyrics/models"
	"text/template"
)

var funcMap = template.FuncMap{
	"add":     func(a, b int) int { return a + b },
	"isAdmin": func(role models.UserRole) bool { return role == models.UserRoleAdmin },
}

const base = "../frontend/src/html/"

var partials = []string{
	base + "partials/template.html",
	base + "partials/navbar.html",
	base + "partials/footer.html",
}

var forumPartials = []string{
	base + "partials/template.html",
	base + "partials/navbar.html",
	base + "partials/footer.html",
	base + "partials/forum_template.html",
	base + "partials/aside.html",
}

type pageConfig struct {
	path     string
	partials []string
}

func parsePage(pagePath string, p []string) (*template.Template, error) {
	files := append([]string{pagePath}, p...)
	return template.New("").Funcs(funcMap).ParseFiles(files...)
}

func ParseTemplates() (map[string]*template.Template, error) {
	pages := map[string]pageConfig{
		"index.html":       {base + "forum/index.html", forumPartials},
		"post.html":        {base + "forum/post/post.html", partials},
		"post-create.html": {base + "forum/post/post-create.html", partials},
		"post-edit.html":   {base + "forum/post/post-edit.html", partials},
		"login.html":       {base + "auth/login.html", partials},
		"register.html":    {base + "auth/register.html", partials},
		"profile.html":     {base + "profile/profile.html", partials},
	}

	templates := make(map[string]*template.Template)
	for name, cfg := range pages {
		t, err := parsePage(cfg.path, cfg.partials)
		if err != nil {
			return nil, err
		}
		templates[name] = t
	}
	return templates, nil
}
