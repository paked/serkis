package serkis

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Public string
}

func (s Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.homeHandler)

	return r
}

func (s Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!\nServing: %s\n", s.Public)
}
