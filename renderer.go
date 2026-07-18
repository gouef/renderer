package renderer

import (
	"log"
	"path/filepath"
	"text/template"

	"github.com/gin-contrib/multitemplate"
	"github.com/gouef/finder"
	"github.com/gouef/renderer/handlers"
	"github.com/gouef/router"
)

type Renderer struct {
	Router          *router.Router
	LayoutPattern   []string
	TemplateDir     string
	TemplateHandler *handlers.TemplateHandler
}

var renderFiles = make([]File, 0)
var customFuncs template.FuncMap

// NewRenderer register and set HTMLRenderer to gouef/router
// Example:
//
//	NewRenderer("./views/templates", []string{"layout", "base.gohtml"})
func NewRenderer(templatesDir string, layoutPattern []string) Renderer {
	if len(layoutPattern) == 0 {
		layoutPattern = []string{"@layout.gohtml", "base.gohtml", "layout.gohtml"}
	}
	return Renderer{TemplateDir: templatesDir, LayoutPattern: layoutPattern}
}

// AddCustomFunc adds a custom function to the template function map
// Example:
//
//	renderer.AddCustomFunc("myFunc", func() string { return "Hello, World!" })
func AddCustomFunc(name string, fn interface{}) {
	if customFuncs == nil {
		customFuncs = make(template.FuncMap)
	}
	customFuncs[name] = fn
}

// AddCustomFuncMap adds a map of custom functions to the template function map
// Example:
//
//	renderer.AddCustomFuncMap(map[string]interface{}{
//		"myFunc": func() string { return "Hello, World!" },
//	})
func AddCustomFuncMap(funcMap map[string]interface{}) {
	for name, fn := range funcMap {
		AddCustomFunc(name, fn)
	}
}

// RegisterRouter register and set HTMLRenderer to gouef/router
// Example:
//
//	renderer.RegisterRouter(r)
func (renderer Renderer) RegisterRouter(r *router.Router) Renderer {
	renderer.Router = r
	templateHandler := &handlers.TemplateHandler{Router: r}
	templateHandler.Initialize()
	renderer.TemplateHandler = templateHandler

	r.SetHtmlRenderer(renderer.HtmlRenderer())
	return renderer
}

func (renderer Renderer) HtmlRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	templatesDir := renderer.TemplateDir
	templateHandler := renderer.TemplateHandler

	for name, fn := range customFuncs {
		templateHandler.AddCustomFunc(name, fn)
	}
	funcMap := templateHandler.GetFuncMap()
	tmpDir := filepath.Join(filepath.Dir(templatesDir), filepath.Base(templatesDir))

	find := finder.FindFiles("*.gohtml").In(templatesDir)

	includes := map[string]*finder.Info{}

	includes = find.Exclude(renderer.LayoutPattern...).Get()

	for p, l := range includes {
		layout, err := findLayout2(l, templatesDir, renderer)

		if err != nil {
			log.Println(err.Error())
		} else {
			f := File{Path: p, Layout: layout}

			f.Includes = findRelevantIncludes(f, templatesDir, renderer)
			renderFiles = append(renderFiles, f)
		}
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range renderFiles {
		l, _ := filepath.Rel(tmpDir, include.Path)

		relatives := make([]string, 0)

		relatives = append(relatives, include.Layout)

		for _, r := range include.Includes {
			relatives = append(relatives, r.Path)
		}

		r.AddFromFilesFuncs(l, funcMap, relatives...)
	}
	return r
}
