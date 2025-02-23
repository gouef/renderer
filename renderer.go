package renderer

import (
	"github.com/gin-contrib/multitemplate"
	"github.com/gouef/renderer/handlers"
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

func loadTemplates(templatesDir string, templateHandler *handlers.TemplateHandler) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcMap := templateHandler.GetFuncMap()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*.gohtml")
	if err != nil {
		panic(err.Error())
	}
	rootLayouts, err := filepath.Glob(templatesDir + "/layout.gohtml")
	if err != nil {
		panic(err.Error())
	}
	layouts = append(layouts, rootLayouts...)

	includes, err := filepath.Glob(templatesDir + "/includes/*.gohtml")
	if err != nil {
		panic(err.Error())
	}

	//	includesNamed, err := filepath.Glob(templatesDir + "/**/*.gohtml")
	/*
		if err != nil {
			panic(err.Error())
		}
		includes = append(includes, includesNamed...)
	*/
	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, includes...)
		r.AddFromFilesFuncs(filepath.Base(include), funcMap, files...)
	}
	return r
}
