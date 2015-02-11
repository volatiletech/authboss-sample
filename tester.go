package main

import (
	"net/http"
	"os"

	"log"

	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/remember"

	"html/template"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type SessionFlasher struct{}

func main() {
	c := authboss.NewConfig()
	cookieStore = securecookie.New([]byte("very-secret"), nil)
	sessionStore = sessions.NewCookieStore([]byte("asdf"))
	c.Storer = NewMemStorer()
	c.LogWriter = os.Stdout
	c.AuthLoginSuccessRoute = "/dashboard"
	c.CookieStoreMaker = NewCookieStorer
	c.SessionStoreMaker = NewSessionStorer
	c.Mailer = authboss.LogMailer(os.Stdout)

	if err := authboss.Init(c); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", authboss.NewRouter())

	templates, _ := template.ParseFiles("views/dashboard.tpl")
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		sstorer := NewSessionStorer(w, r)

		username, ok := sstorer.Get(authboss.SessionKey)

		data := struct {
			Username   string
			IsLoggedIn bool
		}{username, ok}

		templates.ExecuteTemplate(w, "dashboard.tpl", data)
	})

	http.ListenAndServe("localhost:8080", mux)
}
