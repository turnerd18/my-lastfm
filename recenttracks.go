package main

import (
	"fmt"
	lastfm "github.com/turnerd18/go-lastfm"
	"net/http"
	"strconv"
	"time"
)

func recenttracks(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if ((session.Values["username"] == nil) || (session.Values["username"].(string) == "")) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	api, err := newAPI(r)
	if err != nil {
		fmt.Fprintln(w, "error making new api")
		return
	}
	from, _ := time.Parse(time.RFC822, "01 Jan 01 00:00 CST")
	to := time.Now()
	r.ParseForm()
	var page int
	if ((r.PostFormValue("nextpage") != "") || (r.PostFormValue("prevpage") != "")) {
		page, _ = strconv.Atoi(r.PostFormValue("page"))
	} else {
		page = 1
	}
	username := session.Values["username"].(string)
	tracks, err := api.UserGetRecentTracks(username, 50, page, from.Unix(), to.Unix())
	if err != nil {
		fmt.Fprintln(w, "error getting recent tracks")
		return
	}
	type RecentTracksModel struct {
		Tracks []lastfm.APITrack
		Username string
	}
	t := makeTemplate("recenttracks")
	m := RecentTracksModel{Username: username, Tracks: tracks}
	t.Execute(w, m)
}
