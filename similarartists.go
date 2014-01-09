package main

import (
	"fmt"
	"net/http"

	lastfm "github.com/turnerd18/go-lastfm"
)

func similarartists(w http.ResponseWriter, r *http.Request) {
	api, err := newAPI(r)
	if err != nil {
		fmt.Fprintln(w, "error making new api")
		return
		// error page of some sort
	}

	var artist = "Avicii"
	artists, err := api.ArtistGetSimilar(artist, "", 1, 20)
	if err != nil {
		fmt.Fprintln(w, "error getting similar artists")
		return
	}

	type Model struct {
		Artist string
			Artists []lastfm.APIArtist
			Username string
	}
	m := Model{Artist: artist, Artists: artists}
	session, _ := store.Get(r, "session")
	if session.Values["username"] != nil {
		m.Username = session.Values["username"].(string)
	}
	t := makeTemplate("similarartists")
	t.Execute(w, m)
}
