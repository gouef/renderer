package handlers

import "html/template"

var (
	inSnippet bool
)

func snippetStart(name string) template.HTML {
	if inSnippet {
		panic("Error: snippet already opened.")
		return template.HTML("<p>Error: snippet already opened.</p>")
	}

	inSnippet = true
	return template.HTML("<div id=\"snippet-" + name + "\">")
}

func snippetEnd() template.HTML {
	if !inSnippet {
		panic("Error: snippet not started.")
		return template.HTML("<p>Error: snippet not started.</p>")
	}
	inSnippet = false
	return template.HTML("</div>")
}
