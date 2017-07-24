package serkis

import "html/template"

var (
	editTemplate = genTemplate("edit", editTemplateHTML)
	showTemplate = genTemplate("show", showTemplateHTML)
)

const editTemplateHTML = `
<form method="POST" action="/edit/{{ .Fpath }}">
	<textarea cols="80" rows="30" name="contents">{{ .Fcontents }}</textarea>

	<br>
	<br>

	<input name="message" type="text" value="Updating file..."/>

	<br>
	<br>

	<input type="submit" value="Go!" />
</form>
`

const showTemplateHTML = `
<a href="/edit/{{ .Fpath }}">Edit this file</a>

<br>

{{ .UnescapedFcontents }}
`

type TemplateContents struct {
	Fpath     string
	Fcontents string
}

func (tc TemplateContents) UnescapedFcontents() template.HTML {
	return template.HTML(tc.Fcontents)
}

func genTemplate(name, html string) *template.Template {
	return template.Must(template.New(name).Parse(html))
}