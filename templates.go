package serkis

import (
	"html/template"
	"path"
)

var (
	editTemplate    = genTemplate("edit", editTemplateHTML)
	showTemplate    = genTemplate("show", showTemplateHTML)
	rawShowTemplate = genTemplate("rawShow", rawShowTemplateHTML)
	newTemplate     = genTemplate("new", newTemplateHTML)
)

const style = `
<link href='https://fonts.googleapis.com/css?family=Open+Sans:400,600,600italic,300,300italic,400italic,700,700italic,800,800italic' rel='stylesheet' type='text/css'>

<style>
body {
	font-family: 'Open Sans', sans-serif;
	font-weight: 400;
	font-size: 14px;

	color: #221917;

	background-color: white;
}

@media (min-width:1025px) {
	body {
		width: 50%;

		margin-left: auto;
		margin-right: auto;
	}
}

.links {
	margin-right: 20px;
}
</style>
`

const editTemplateHTML = style + `
<a href="/{{ .Fpath }}" class="links">View this file</a>

<form method="POST" action="/edit/{{ .Fpath }}">
	<textarea cols="80" rows="30" name="contents" accept-charset="UTF-8">{{ .Fcontents }}</textarea>

	<br>
	<br>

	<input name="message" type="text" value="Updating file..."/>

	<br>
	<br>

	<input type="submit" value="Commit" />
</form>
`

const showTemplateHTML = style + `
<a href="/edit/{{ .Fpath }}" class="links">Edit this file</a>
<a href="/new" class="links">Create a new file</a>
<a href="/{{ .BackURL }}" class="links">Back</a>
<a href="/{{ .ShareURL }}" class="links">Share</a>

<br>

{{ .UnescapedFcontents }}
`

const rawShowTemplateHTML = style + `
{{ .UnescapedFcontents }}
`

const newTemplateHTML = style + `
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

func (tc TemplateContents) ShareURL() template.HTML {
	s := Share{Fpath: tc.Fpath}

	secret, _ := s.Secret(cryptoKey)

	return template.HTML("share/" + secret)
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
