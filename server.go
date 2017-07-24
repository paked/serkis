package serkis

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

type HandlerWithFile func(http.ResponseWriter, *http.Request, []byte)

type Server struct {
	Public string
}

func (s Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc(
		"/edit/{rest:.*}",
		s.middlewareGetFile(s.handleShowEdit),
	).Methods("GET")

	r.HandleFunc(
		"/edit/{rest:.*}",
		s.handleEdit,
	).Methods("POST")

	r.HandleFunc(
		"/{rest:.*}",
		s.middlewareGetFile(s.handleShow),
	).Methods("GET")

	return r
}

func (s Server) middlewareGetFile(f HandlerWithFile) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fpath := mux.Vars(req)["rest"]
		w.Header().Set("Content-Type", "text/html")

		raw, err := s.file(fpath)
		if err != nil {
			fmt.Fprintln(w, "Could not find file")
			return
		}

		f(w, req, raw)
	}
}

func (s Server) handleShowEdit(w http.ResponseWriter, req *http.Request, raw []byte) {
	fpath := mux.Vars(req)["rest"]

	t, err := template.New("edit").Parse(editTemplate)
	if err != nil {
		fmt.Fprintln(w, "Could not parse template")
		return
	}

	data := struct {
		Fpath     string
		Fcontents string
	}{
		Fpath:     fpath,
		Fcontents: string(raw),
	}

	err = t.Execute(w, data)

	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}
}

func (s Server) handleEdit(w http.ResponseWriter, req *http.Request) {
	rawfpath := mux.Vars(req)["rest"]
	fpath := s.path(rawfpath)

	contents := req.FormValue("contents")

	err := ioutil.WriteFile(fpath, []byte(contents), 0644)
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}

	http.Redirect(w, req, "/"+rawfpath, 301)
}

func (s Server) handleShow(w http.ResponseWriter, req *http.Request, raw []byte) {
	md := blackfriday.MarkdownCommon(raw)

	fmt.Fprintf(w, "%s", md)
}

func (s Server) file(fpath string) ([]byte, error) {
	return ioutil.ReadFile(s.path(fpath))
}

func (s Server) path(fpath string) string {
	// We calls like `/` and `/streams` to resolve to `/README.md` and `/streams/README.md` respectively.
	if path.Ext(fpath) != ".md" {
		fpath = path.Join("README.md")
	}

	return path.Join(s.Public, path.Clean(fpath))
}
