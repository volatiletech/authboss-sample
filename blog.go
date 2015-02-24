package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	ab "gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/register"
	_ "gopkg.in/authboss.v0/remember"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/unrolled/render"
)

func setupAuthboss() {
	ab.Cfg.Storer = NewMemStorer()
	ab.Cfg.MountPath = "/auth"
	ab.Cfg.LogWriter = os.Stdout

	ab.Cfg.AuthLoginSuccessRoute = "/"
	ab.Cfg.Layout = template.Must(template.New("layout").ParseFiles("views/layout.tpl"))

	ab.Cfg.CookieStoreMaker = NewCookieStorer
	ab.Cfg.SessionStoreMaker = NewSessionStorer

	ab.Cfg.Mailer = ab.LogMailer(os.Stdout)

	ab.Cfg.XSRFName = "csrf_token"
	ab.Cfg.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

var rendering = render.New(render.Options{
	Directory:     "views",
	Layout:        "layout",
	Extensions:    []string{".tpl"},
	IsDevelopment: "true",
})

func main() {
	// Initialize Sessions and Cookies
	cookieStore = securecookie.New([]byte("very-secret"), nil)
	sessionStore = sessions.NewCookieStore([]byte("asdf"))

	// Initialize ab.
	setupab()

	// Create our router and middleware chain
	mux := mux.NewRouter()
	stack := alice.New(logger, nosurfHandler).Then(mux)

	// Routes
	gets := mux.Methods("GET").Subrouter()
	posts := mux.Methods("POST").Subrouter()

	mux.Handle("/auth", ab.NewRouter())

	gets.HandleFunc("/blogs/:id/edit", edit)
	gets.HandleFunc("/blogs/:id", show)
	gets.HandleFunc("/blogs", index)

	posts.HandleFunc("/blogs/:id", update)
	posts.HandleFunc("/blogs", create)
	mux.Methods("DELETE").HandleFunc("/blogs/:id", destroy)

	gets.HandleFunc("/", index)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func show(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func new(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func create(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func edit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func update(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}

func destroy(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
}
