package serkis

const editTemplate = `
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
