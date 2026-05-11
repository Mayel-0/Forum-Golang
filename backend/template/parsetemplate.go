package template

import (
	"fmt"
	"lyrics/models"
	"text/template"
	"time"
)

var funcMap = template.FuncMap{
	"add":     func(a, b int) int { return a + b },
	"isAdmin": func(role models.UserRole) bool { return role == models.UserRoleAdmin },
	"timeAgo": func(t time.Time) string {
		d := time.Since(t)
		switch {
		case d < time.Minute:
			return "à l'instant"
		case d < time.Hour:
			m := int(d.Minutes())
			if m == 1 {
				return "il y a 1 min"
			}
			return fmt.Sprintf("il y a %d min", m)
		case d < 24*time.Hour:
			h := int(d.Hours())
			if h == 1 {
				return "il y a 1 h"
			}
			return fmt.Sprintf("il y a %d h", h)
		case d < 7*24*time.Hour:
			days := int(d.Hours() / 24)
			if days == 1 {
				return "hier"
			}
			return fmt.Sprintf("il y a %d j", days)
		default:
			return t.Format("02/01/2006")
		}
	},
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
		"post.html":        {base + "forum/post/post.html", forumPartials},
		"subcategory.html": {base + "forum/post/subcategory.html", forumPartials},
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
