package handlers

import (
	"github.com/gouef/router"
	"html/template"
	"log"
)

type TemplateHandler struct {
	Router *router.Router
}

type UrlForFunc func(name string, params ...interface{}) string

func (t *TemplateHandler) Initialize() {
	n := t.Router.GetNativeRouter()
	n.SetFuncMap(t.GetFuncMap())
}

func (t *TemplateHandler) GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"snippet":    snippetStart,
		"snippetEnd": snippetEnd,
		"endSnippet": snippetEnd,
		"link": UrlForFunc(func(name string, params ...interface{}) string {
			paramMap := make(map[string]interface{})
			if len(params) > 0 {
				for i := 0; i < len(params); i += 2 {
					key := params[i].(string)
					value := params[i+1]
					paramMap[key] = value
				}
			}

			url, err := t.Router.GenerateUrlByName(name, paramMap)
			if err != nil {
				log.Println("Error generating URL for", name, ":", err)
				return "/error"
			}
			return url
		}),
	}
}
