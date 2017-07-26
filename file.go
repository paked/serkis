package serkis

import (
	"io/ioutil"
	"os"
	"path"
)

type File struct {
	Public string
	Path   string
}

// RedirectTo detects if a file needs to be redirected somewhere, and handles it.
func (f File) RedirectTo() (string, bool, error) {
	fi, err := os.Stat(f.LocalPath())

	switch {
	case err != nil:
		return "", false, err
	case fi.IsDir():
		return path.Join(f.VPath(), "README.md"), true, nil
	default:
		return "", false, nil
	}
}

// Contents gets the contents of a file.
func (f File) Contents() ([]byte, error) {
	return ioutil.ReadFile(f.LocalPath())
}

// VPath gets the virtual path of a file, as it would be represented on a route.
func (f File) VPath() string {
	return path.Clean(f.Path)
}

// LocalPath gets the path of the file in relation to the file system
func (f File) LocalPath() string {
	return path.Join(f.Public, f.VPath())
}

// TemplateData is a helper method which provides a quick way to get data for templates.
func (f File) TemplateData() TemplateContents {
	contents, _ := f.Contents()

	return TemplateContents{
		Fpath:     f.VPath(),
		Fcontents: string(contents),
	}
}
