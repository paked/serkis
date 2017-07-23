package main

import (
	"net/http"

	"github.com/paked/configure"
	"github.com/paked/serkis"
)

var (
	conf = configure.New()

	public = conf.String("public", "public", "The files to serve")
)

func init() {
	conf.Use(
		configure.NewFlag(),
		configure.NewEnvironment(),
	)
}

func main() {
	conf.Parse()

	s := serkis.Server{
		Public: *public,
	}

	http.ListenAndServe("0.0.0.0:8765", s.Router())
}
