package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	ab "gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/register"
	_ "gopkg.in/authboss.v0/remember"

	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

func setupAuthboss() {
	ab.Cfg.Storer = NewMemStorer()
	ab.Cfg.MountPath = "/auth"
	ab.Cfg.LogWriter = os.Stdout

	ab.Cfg.AuthLoginSuccessRoute = "/"
	layout := template.Must(template.New("layout").
		Funcs(funcs).
		ParseFiles("views/layout.html.tpl"))
	ab.Cfg.Layout = layout

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

var (
	templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	schemaDec = schema.NewDecoder()
)

func main() {
	// Initialize Sessions and Cookies
	cookieStore = securecookie.New([]byte("very-secret"), nil)
	sessionStore = sessions.NewCookieStore([]byte("asdf"))

	// Initialize ab.
	setupAuthboss()

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)
	mux := mux.NewRouter()

	// Routes
	gets := mux.Methods("GET").Subrouter()
	posts := mux.Methods("POST").Subrouter()

	mux.Handle("/auth", ab.NewRouter())

	gets.HandleFunc("/blogs/new", new)
	gets.HandleFunc("/blogs/{id}/edit", edit)
	gets.HandleFunc("/blogs/{id}", show)
	gets.HandleFunc("/blogs", index)
	gets.HandleFunc("/", index)

	posts.HandleFunc("/blogs/{id}", update)
	posts.HandleFunc("/blogs", create)

	// This should actually be a destroys.X but I can't be bothered to make a proper
	// destroy link using javascript atm.
	gets.HandleFunc("/blogs/{id}/destroy", destroy)

	mux.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, "Not found")
	})

	// Set up our middleware chain
	stack := alice.New(logger /*, nosurfing*/).Then(mux)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}

func layoutData(w http.ResponseWriter, r *http.Request) ab.HTMLData {
	return ab.HTMLData{
		"loggedin":   false,
		"username":   "",
		"csrf_token": nosurf.Token(r),
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, "index", data)
}

func show(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	data := layoutData(w, r).MergeKV("post", blogs[id], "id", id)
	err := templates.Render(w, "show", data)
	if err != nil {
		panic(err)
	}
}

func new(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("post", Blog{})
	mustRender(w, "new", data)
}

func create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	// TODO: Validation

	var b Blog
	if badRequest(w, schemaDec.Decode(&b, r.PostForm)) {
		return
	}

	b.Date = time.Now()
	b.AuthorID = "Zeratul"

	blogs = append(blogs, b)

	data := layoutData(w, r).MergeKV("post", b)
	mustRender(w, "show", data)
}

func edit(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	data := layoutData(w, r).MergeKV("post", blogs[id], "id", id)
	mustRender(w, "edit", data)
}

func update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	id, ok := blogID(w, r)
	if !ok {
		return
	}

	// TODO: Validation

	var b = blogs[id]
	if badRequest(w, schemaDec.Decode(&b, r.PostForm)) {
		return
	}

	b.Date = time.Now()

	blogs[id] = b

	data := layoutData(w, r).MergeKV("post", blogs[id], "id", id)
	mustRender(w, "show", data)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	if len(blogs) == 1 {
		blogs = []Blog{}
	} else {
		for i := id; i < len(blogs)-1; i++ {
			blogs[i], blogs[i+1] = blogs[i+1], blogs[i]
		}
		blogs = blogs[:len(blogs)-1]
	}

	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, "index", data)
}

func blogID(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	str := vars["id"]

	id, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error parsing blog id:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	if id < 0 || id >= len(blogs) {
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	return id, true
}

func mustRender(w http.ResponseWriter, name string, data interface{}) {
	err := templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}
