package serkis

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

type HandlerWithFile func(http.ResponseWriter, *http.Request, File)

type Server struct {
	Public string

	HTTPUsername string
	HTTPPassword string
	HTTPRealm    string

	Git *Git

	GitHubWebhookSecret string
}

func (s Server) Init() error {
	err := s.Git.Clone(s.Public)
	if err != nil {
		return err
	}

	err = s.Git.Config(s.Public, "user.name", s.Git.AuthorName)
	if err != nil {
		return err
	}

	err = s.Git.Config(s.Public, "user.email", s.Git.AuthorEmail)
	if err != nil {
		return err
	}

	return nil
}

func (s Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/gh_webhook", s.handleWebhook)
	r.HandleFunc(
		"/share/{secret:.*}",
		s.handleSecret,
	).Methods("GET")

	r.HandleFunc(
		"/new",
		s.middlewareBasicAuth(s.handleShowNew),
	).Methods("GET")

	r.HandleFunc(
		"/new",
		s.middlewareBasicAuth(s.handleNew),
	).Methods("POST")

	r.HandleFunc(
		"/edit/{rest:.*}",
		s.middlewareBasicAuth(s.middlewareGetFile(s.handleShowEdit)),
	).Methods("GET")

	r.HandleFunc(
		"/edit/{rest:.*}",
		s.middlewareBasicAuth(s.handleEdit),
	).Methods("POST")

	r.HandleFunc(
		"/{rest:.*}",
		s.middlewareBasicAuth(s.middlewareGetFile(s.handleShow)),
	).Methods("GET")

	return r
}

func (s Server) middlewareBasicAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", s.HTTPRealm))

			http.Error(w, "Could not authorize user", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(auth[1])
		if err != nil {
			fmt.Fprintln(w, "Could not parse auth header:", err)
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !s.authed(pair[0], pair[1]) {
			http.Error(w, "Could not authorize user", http.StatusUnauthorized)
			return
		}

		f(w, req)
	}
}

func (s Server) authed(username, password string) bool {
	return username == s.HTTPUsername && password == s.HTTPPassword
}

func (s Server) middlewareGetFile(f HandlerWithFile) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		fi := File{
			Path:   mux.Vars(req)["rest"],
			Public: s.Public,
		}

		n, should, err := fi.RedirectTo()
		if err != nil {
			fmt.Fprintln(w, "Error getting file:", err)
			return
		}

		if should {
			// This redirect will not properly handle redirects to the edit page.
			http.Redirect(w, req, n, 301)
			return
		}

		f(w, req, fi)
	}
}

func (s Server) handleSecret(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	secret := mux.Vars(req)["secret"]

	sh, err := NewShareFromSecret(secret, cryptoKey)
	if err != nil {
		fmt.Fprintln(w, "Failed to decrypt secret:", err)
		return
	}

	fi := File{
		Path:   sh.Fpath,
		Public: s.Public,
	}

	_, _, err = fi.RedirectTo()
	if err != nil {
		fmt.Fprintln(w, "Error getting file:", err)
		return
	}

	data := fi.TemplateData()

	data.Fcontents = string(blackfriday.MarkdownCommon([]byte(data.Fcontents)))

	err = rawShowTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}
}

func (s Server) handleShowEdit(w http.ResponseWriter, req *http.Request, fi File) {
	err := editTemplate.Execute(w, fi.TemplateData())
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}
}

func (s Server) handleShowNew(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := newTemplate.Execute(w, nil)
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}
}

func (s Server) handleNew(w http.ResponseWriter, req *http.Request) {
	fi := File{
		Path:   req.FormValue("path"),
		Public: s.Public,
	}

	err := ioutil.WriteFile(fi.LocalPath(), []byte{}, 0644)
	if err != nil {
		fmt.Fprintln(w, "Could not create new file:", err)
		return
	}

	http.Redirect(w, req, path.Join("edit", fi.VPath()), 301)
}

func (s Server) handleEdit(w http.ResponseWriter, req *http.Request) {
	fi := File{
		Path:   mux.Vars(req)["rest"],
		Public: s.Public,
	}

	r := strings.NewReplacer("\r\n", "\n")

	contents := r.Replace(req.FormValue("contents"))
	message := req.FormValue("message")

	err := ioutil.WriteFile(fi.LocalPath(), []byte(contents), 0644)
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}

	go s.Git.PushNewChanges(s.Public, fi.VPath(), message)

	http.Redirect(w, req, path.Join("/", fi.VPath()), 301)
}

func (s Server) handleShow(w http.ResponseWriter, req *http.Request, fi File) {
	data := fi.TemplateData()

	data.Fcontents = string(blackfriday.MarkdownCommon([]byte(data.Fcontents)))

	err := showTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, "Failed to render template: ", err)
		return
	}
}

func (s Server) handleWebhook(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "application/json")

	wv := WebhookValidator{
		Secret:    s.GitHubWebhookSecret,
		Signature: req.Header.Get("x-hub-signature"),
		Body:      req.Body,
	}

	ok, err := wv.Valid()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid webhook: ", err)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Incorrect secret")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "{}")

	go s.Git.PullRemoteChanges(s.Public)
}
