package main

import (
	"net/http"

	"github.com/paked/configure"
	"github.com/paked/serkis"
)

var (
	conf = configure.New()

	public = conf.String("public", "static", "The files to serve")
	port   = conf.String("port", "8765", "The port to serve from")

	httpUsername = conf.String("http-username", "admin", "Username for HTTP authentication")
	httpPassword = conf.String("http-password", "admin", "Password for HTTP basic authentication")
	httpRealm    = conf.String("http-realm", "serkis", "Realm for HTTP basic authentication")

	gitURL         = conf.String("git-url", "", "URL of git repository")
	gitUsername    = conf.String("git-username", "", "Username for git account")
	gitPassword    = conf.String("git-password", "", "Password for git account")
	gitAuthorName  = conf.String("git-author-name", "", "Name for git account")
	gitAuthorEmail = conf.String("git-author-email", "", "Email for git account")

	githubWebhookSecret = conf.String("github-webhook-secret", "", "Secret for GitHub Webhook")
)

func init() {
	conf.Use(
		configure.NewFlag(),
		configure.NewEnvironment(),
		configure.NewJSONFromFile("config.json"),
	)
}

func main() {
	conf.Parse()

	s := serkis.Server{
		Public: *public,

		HTTPUsername: *httpUsername,
		HTTPPassword: *httpPassword,
		HTTPRealm:    *httpRealm,

		Git: &serkis.Git{
			URL:         *gitURL,
			Username:    *gitUsername,
			Password:    *gitPassword,
			AuthorName:  *gitAuthorName,
			AuthorEmail: *gitAuthorEmail,
		},

		GitHubWebhookSecret: *githubWebhookSecret,
	}

	s.Init()

	http.ListenAndServe("0.0.0.0:"+*port, s.Router())
}
