package serkis

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

type Server struct {
	Public string
}

func (s Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/edit/{rest:.*}", s.handleEdit).Methods("GET", "PUT")
	r.HandleFunc("/{rest:.*}", s.handleShow).Methods("GET")

	return r
}

func (s Server) handleEdit(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Editing: "+mux.Vars(req)["rest"])
}

func (s Server) handleShow(w http.ResponseWriter, req *http.Request) {
	fpath := mux.Vars(req)["rest"]

	f, err := s.file(fpath)

	if err != nil {
		fmt.Fprintln(w, "Could not find file")
		return
	}

	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, "%s", f)
}

func (s Server) file(fpath string) ([]byte, error) {
	p := path.Join(s.Public, path.Clean(fpath))

	raw, err := ioutil.ReadFile(p)
	if err != nil {
		return []byte{}, err
	}

	md := blackfriday.MarkdownCommon(raw)

	return md, nil
}
