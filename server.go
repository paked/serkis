package serkis

import (
	"fmt"
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

	r.HandleFunc("/edit/{rest:.*}", s.handleEdit).Methods("GET", "PUT")

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

func (s Server) handleEdit(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Editing: "+mux.Vars(req)["rest"])
}

func (s Server) handleShow(w http.ResponseWriter, req *http.Request, raw []byte) {
	md := blackfriday.MarkdownCommon(raw)

	fmt.Fprintf(w, "%s", md)
}

func (s Server) file(fpath string) ([]byte, error) {
	p := path.Join(s.Public, path.Clean(fpath))

	return ioutil.ReadFile(p)
}
