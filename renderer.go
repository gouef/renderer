package renderer

import (
	"errors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gouef/finder"
	"github.com/gouef/renderer/handlers"
	"github.com/gouef/router"
	"log"
	"path/filepath"
)

type File struct {
	Path     string
	Layout   string
	Includes []File
}

var renderFiles = make([]File, 0)

// RegisterToRouter register and set HTMLRenderer to gouef/router
// Example:
//
//	RegisterToRouter(r, "./views/templates")
func RegisterToRouter(r *router.Router, templatesDir string) {
	templateHandler := &handlers.TemplateHandler{Router: r}
	templateHandler.Initialize()
	r.SetHtmlRenderer(LoadTemplates(templatesDir, templateHandler))
}

// LoadTemplates
// Example:
//
//	r := router.NewRouter()
//	templateHandler := &handlers.TemplateHandler{Router: r}
//	templateHandler.Initialize()
//	r.SetHtmlRenderer(renderer.LoadTemplates("./views/templates", templateHandler))
func LoadTemplates(templatesDir string, templateHandler *handlers.TemplateHandler) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcMap := templateHandler.GetFuncMap()
	tmpDir := filepath.Join(filepath.Dir(templatesDir), filepath.Base(templatesDir))

	find := finder.FindFiles("*.gohtml").In(templatesDir)

	includes := map[string]*finder.Info{}

	includes = find.Exclude("layout.gohtml", "@layout.gohtml", "base.gohtml").Get()

	for p, l := range includes {
		layout, err := findLayout2(l, templatesDir)

		if err != nil {
			log.Println(err.Error())
		} else {
			f := File{Path: p, Layout: layout}

			f.Includes = findRelevantIncludes(f, templatesDir)
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

func findRelevantIncludes(file File, searchDir string) (result []File) {
	dir := filepath.Dir(file.Path)
	absDir, _ := filepath.Abs(dir)
	files := findIncludesInDir(absDir, searchDir)

	for p := range files {
		result = append(result, File{Path: p})
	}

	return result
}

func findIncludesInDir(dir string, templateDir string) map[string]*finder.Info {
	templateDir, _ = filepath.Abs(templateDir)
	find := finder.In(dir).FindFiles("*.gohtml").NotRecursive()
	files := find.Exclude("layout.gohtml", "@layout.gohtml", "base.gohtml").Get()

	if dir != templateDir {
		for p, i := range findIncludesInDir(filepath.Dir(dir), templateDir) {
			files[p] = i
		}
	}

	return files
}

func findLayout2(file *finder.Info, templatesDir string) (string, error) {
	dir := filepath.Dir(file.Path)
	absDir, _ := filepath.Abs(dir)
	absTemplate, _ := filepath.Abs(templatesDir)
	return findLayoutInDir(absDir, absTemplate)
}

func findLayoutInDir(dir string, templateDir string) (string, error) {
	find := finder.In(dir).FindFiles("*.gohtml").NotRecursive()
	files := find.Match("layout.gohtml", "base.gohtml")
	if len(files) >= 1 {
		first := FirstRecord(files)
		return first, nil
	}

	if dir != templateDir {
		return findLayoutInDir(filepath.Dir(dir), templateDir)
	}

	return "", errors.New("not found layout")
}

func FirstRecord(m map[string]*finder.Info) string {
	for k := range m {
		return k
	}
	return ""
}
