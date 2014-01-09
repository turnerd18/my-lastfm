package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
)

func login(w http.ResponseWriter, r *http.Request) {
	// if login form was submitted
	r.ParseForm()
	if((r.PostFormValue("username") != "") && (r.PostFormValue("password") != "")) {
		// return errors as form validation
		api, err := newAPI(r)
		if err != nil {
			fmt.Fprintln(w, "error making new api")
			return
		}
		err = api.AuthGetMobileSession(r.PostFormValue("username"), r.PostFormValue("password"))
		if err != nil {
			fmt.Fprintln(w, "error logging into last.fm")
			return
		}
		// check if user in database
		con, err := sql.Open("mysql", dbstring)
		if err != nil {
			con.Close()
			return
		}
		var skdb string
		err = con.QueryRow("SELECT sk FROM users WHERE user=?", api.Username).Scan(&skdb)
		if err != nil {
			// insert user into database
			_, err = con.Exec("INSERT INTO users (user, sk) VALUES (?, ?)", api.Username, api.Sk)
			if err != nil {
				fmt.Fprintln(w, "error inserting new user record: " + err.Error())
				return
			}
		} else {
			// if session key has changed, update database
			if api.Sk != skdb {
				_, err = con.Exec("UPDATE users SET sk=? WHERE user=?", api.Sk, api.Username)
				if err != nil {
					fmt.Fprintln(w, "error updating session key: " + err.Error())
					return
				}
			}
		}
		session, _ := store.Get(r, "session")
		session.Values["username"] = api.Username
		session.Values["sk"] = api.Sk
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		// load the login form
		t := makeTemplate("login")
		t.Execute(w, nil)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = ""
	session.Values["sk"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}
