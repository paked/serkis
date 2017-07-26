package serkis

import (
	"html/template"
	"path"
)

var (
	editTemplate = genTemplate("edit", editTemplateHTML)
	showTemplate = genTemplate("show", showTemplateHTML)
	newTemplate  = genTemplate("new", newTemplateHTML)
)

const editTemplateHTML = `
<a href="/{{ .Fpath }}">View this file</a>

<form method="POST" action="/edit/{{ .Fpath }}">
	<textarea cols="80" rows="30" name="contents">{{ .Fcontents }}</textarea>

	<br>
	<br>

	<input name="message" type="text" value="Updating file..."/>

	<br>
	<br>

	<input type="submit" value="Commit" />
</form>
`

const showTemplateHTML = `
<a href="/edit/{{ .Fpath }}">Edit this file</a>
<a href="/new">Create a new file</a>
<a href="/{{ .BackURL }}">Back</a>

<br>

{{ .UnescapedFcontents }}
`

const newTemplateHTML = `
<form method="POST" action="/new">
	<input name="path" type="text"/>

	<br>

	<input type="submit" value="Create file" />
</form>
`

type TemplateContents struct {
	Fpath     string
	Fcontents string
}

func (tc TemplateContents) UnescapedFcontents() template.HTML {
	return template.HTML(tc.Fcontents)
}

func (tc TemplateContents) BackURL() string {
	url := path.Dir(tc.Fpath)

	if path.Base(tc.Fpath) == "README.md" {
		url = path.Dir(url)
	}

	return url
}

func genTemplate(name, html string) *template.Template {
	return template.Must(template.New(name).Parse(html))
}
