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
	"github.com/justinas/nosurf"
)

type SessionFlasher struct{}

func main() {
	cookieStore = securecookie.New([]byte("very-secret"), nil)
	sessionStore = sessions.NewCookieStore([]byte("asdf"))
	authboss.Cfg.Storer = NewMemStorer()
	authboss.Cfg.LogWriter = os.Stdout
	authboss.Cfg.AuthLoginSuccessRoute = "/dashboard"
	authboss.Cfg.CookieStoreMaker = NewCookieStorer
	authboss.Cfg.SessionStoreMaker = NewSessionStorer
	authboss.Cfg.Mailer = authboss.LogMailer(os.Stdout)
	authboss.Cfg.XSRFName = "csrf_token"
	authboss.Cfg.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	if err := authboss.Init(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", nosurf.New(authboss.NewRouter()))

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
