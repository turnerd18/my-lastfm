package main

import (
	"database/sql"
	"fmt"
	"text/template"
	"net/http"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type track struct {
	Name, Artist, Album, Date, Image string
}

func index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if ((session.Values["username"] == nil) || (session.Values["username"].(string) == "")) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	con, err := sql.Open("mysql", dbstring)
	if err != nil {
		con.Close()
		return
	}
	username := session.Values["username"].(string)
	rows, err := con.Query("SELECT name, artist, album, date, image FROM song_discovery WHERE user=? ORDER BY date DESC LIMIT 0,100", username)
	if err != nil {
		fmt.Println("error selecting: " + err.Error())
		return
	}

	var tracks []track
	var t track
	for rows.Next() {
		rows.Scan(&t.Name, &t.Artist, &t.Album, &t.Date, &t.Image)
		timestamp, _ := strconv.ParseInt(t.Date, 10, 64)
		format := "Jan 2, 2006 at 3:04pm"
		t.Date = time.Unix(timestamp, 0).Format(format)
		tracks = append(tracks, t)
	}

	tmpl := makeTemplate("index")
	type Model struct {
		Tracks []track
		Username string
	}
	m := Model{Username: username, Tracks: tracks}
	tmpl.Execute(w, m)
}

func add(a, b int) int {
	return a + b
}

func indexmore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if(r.PostFormValue("page") == "") {
		fmt.Fprintf(w, "ERROR: no page set")
		return
	}
	session, _ := store.Get(r, "session")
	if ((session.Values["username"] == nil) || (session.Values["username"].(string) == "")) {
		fmt.Fprintf(w, "ERROR no username found")
		return
	}
	con, err := sql.Open("mysql", dbstring)
	if err != nil {
		con.Close()
		fmt.Fprintf(w, "ERROR could not connect to database: %q", err.Error())
		return
	}

	username := session.Values["username"].(string)
	limit, _ := strconv.Atoi(r.PostFormValue("page"))
	limit *= 100
	rows, err := con.Query("SELECT name, artist, album, date, image FROM song_discovery WHERE user=? ORDER BY date DESC LIMIT ?,?", username, limit + 1, 100)
	if err != nil {
		fmt.Fprintf(w, "ERROR selecting tracks: %q", err.Error())
		return
	}

	var tracks []track
	var t track
	for rows.Next() {
		rows.Scan(&t.Name, &t.Artist, &t.Album, &t.Date, &t.Image)
		timestamp, _ := strconv.ParseInt(t.Date, 10, 64)
		format := "Jan, 2 2006 at 3:04pm"
		t.Date = time.Unix(timestamp, 0).Format(format)
		tracks = append(tracks, t)
	}

	type Model struct {
		Tracks []track
		Limit int
	}
	model := Model{Tracks: tracks, Limit: limit + 1}
	funcmap := template.FuncMap{"add": func(a, b int) int {
		return a + b
		}}
	tmpl := template.Must(template.New("example").Funcs(funcmap).Parse(indexmoretmpl))
	//TODO Why won't it parse this file correctly????
	//tmpl := template.Must(template.New("example").Funcs(funcmap).ParseFiles("templates/indexmore.html"))
	tmpl.Execute(w, model)
	con.Close()
}

// declaring template as variable because indexmore.html wont parse correctly
var indexmoretmpl = `{{ range .Tracks }}
<li>
<img src="{{ .Image }}" /><br />
{{ .Artist }}<br />
{{ .Name }}<br />
{{ .Date }}<br />
</li>
{{ end }}`
