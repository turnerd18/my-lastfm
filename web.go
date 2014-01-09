package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(sessionsecret))

func handler() {
	http.HandleFunc("/", index)
	http.HandleFunc("/indexmore", indexmore)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/recenttracks", recenttracks)
	http.HandleFunc("/similarartists", similarartists)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
}

func main() {
	handler()
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
