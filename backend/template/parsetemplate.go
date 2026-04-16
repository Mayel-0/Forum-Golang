package template

import "text/template"

var tpl *template.Template
var err error

func ParseTemplates() (*template.Template, error) {
	tpl, err = template.ParseFiles(
		"../frontend/src/html/acceuil.html",
	)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}
