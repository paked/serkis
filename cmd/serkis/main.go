package main

import (
	"net/http"

	"github.com/paked/configure"
	"github.com/paked/serkis"
)

var (
	conf = configure.New()

	public = conf.String("public", "public", "The files to serve")
	port   = conf.String("port", "8765", "The port to serve from")

	httpUsername = conf.String("http-username", "admin", "Username for HTTP authentication")
	httpPassword = conf.String("http-password", "admin", "Password for HTTP basic authentication")
	httpRealm    = conf.String("http-realm", "serkis", "Realm for HTTP basic authentication")
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

		HTTPUsername: *httpUsername,
		HTTPPassword: *httpPassword,
		HTTPRealm:    *httpRealm,
	}

	http.ListenAndServe("0.0.0.0:"+*port, s.Router())
}
