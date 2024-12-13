package renderer

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	layoutPath := filepath.Join("templates", "@layout.gohtml")
	tmplPath := filepath.Join("templates", tmpl+".gohtml")

	t, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
