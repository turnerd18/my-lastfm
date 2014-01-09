package main

import (
	"html/template"
	"net/http"

	lastfm "github.com/turnerd18/go-lastfm"
)

func makeTemplate(content string) *template.Template {
	return template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/" + content + ".html"))
}

func newAPI(r *http.Request) (*lastfm.API, error) {
	var username, sk string
	session, _ := store.Get(r, "session")
	if session.Values["username"] != nil {
		username = session.Values["username"].(string)
	}
	if session.Values["sk"] != nil {
		sk = session.Values["sk"].(string)
	}
	api, err := lastfm.NewAPI(apikey, apisecret, username, sk)
	return api, err
}
