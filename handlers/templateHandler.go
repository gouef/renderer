package handlers

import (
	"html/template"
	"log"

	"github.com/gouef/router"
)

type TemplateHandler struct {
	Router      *router.Router
	customFuncs template.FuncMap
}

type UrlForFunc func(name string, params ...interface{}) string

func (t *TemplateHandler) Initialize() {
	if t.Router != nil {
		n := t.Router.GetNativeRouter()
		n.SetFuncMap(t.GetDefaultFuncMap())
	}
}

func (t *TemplateHandler) AddCustomFunc(name string, fn interface{}) {
	if t.customFuncs == nil {
		t.customFuncs = make(template.FuncMap)
	}
	t.customFuncs[name] = fn
	if t.Router != nil {
		n := t.Router.GetNativeRouter()
		n.SetFuncMap(t.GetFuncMap())
	}
}

func (t *TemplateHandler) GetDefaultFuncMap() template.FuncMap {
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

func (t *TemplateHandler) GetFuncMap() template.FuncMap {
	funcMap := t.GetDefaultFuncMap()

	if t.customFuncs != nil {
		for name, fn := range t.customFuncs {
			funcMap[name] = fn
		}
	}

	return funcMap
}
