package renderer

import (
	"errors"
	"github.com/gouef/finder"
	"path/filepath"
)

func findRelevantIncludes(file File, searchDir string, renderer Renderer) (result []File) {
	dir := filepath.Dir(file.Path)
	absDir, _ := filepath.Abs(dir)
	files := findIncludesInDir(absDir, searchDir, renderer)

	for p := range files {
		result = append(result, File{Path: p})
	}

	return result
}

func findIncludesInDir(dir string, templateDir string, renderer Renderer) map[string]*finder.Info {
	templateDir, _ = filepath.Abs(templateDir)
	find := finder.In(dir).FindFiles("*.gohtml").NotRecursive()
	files := find.Exclude(renderer.LayoutPattern...).Get()

	if dir != templateDir {
		for p, i := range findIncludesInDir(filepath.Dir(dir), templateDir, renderer) {
			files[p] = i
		}
	}

	return files
}

func findLayout2(file *finder.Info, templatesDir string, renderer Renderer) (string, error) {
	dir := filepath.Dir(file.Path)
	absDir, _ := filepath.Abs(dir)
	absTemplate, _ := filepath.Abs(templatesDir)
	return findLayoutInDir(absDir, absTemplate, renderer)
}

func findLayoutInDir(dir string, templateDir string, renderer Renderer) (string, error) {
	find := finder.In(dir).FindFiles("*.gohtml").NotRecursive()
	files := find.Match(renderer.LayoutPattern...)
	if len(files) >= 1 {
		first := FirstRecord(files)
		return first, nil
	}

	if dir != templateDir {
		return findLayoutInDir(filepath.Dir(dir), templateDir, renderer)
	}

	return "", errors.New("not found layout")
}

func FirstRecord(m map[string]*finder.Info) string {
	for k := range m {
		return k
	}
	return ""
}
