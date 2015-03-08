package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	ab "gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/confirm"
	_ "gopkg.in/authboss.v0/lock"
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

var (
	database  = NewMemStorer()
	templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	schemaDec = schema.NewDecoder()
)

func setupAuthboss() {
	ab.Cfg = ab.NewConfig()
	ab.Cfg.Storer = database
	ab.Cfg.MountPath = "/auth"
	ab.Cfg.LogWriter = os.Stdout
	ab.Cfg.HostName = `http://localhost:3000`

	ab.Cfg.LayoutDataMaker = layoutData

	b, err := ioutil.ReadFile(filepath.Join("views", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Cfg.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b)))
	//ab.Cfg.LayoutEmail = template.Must(template.New("layout").Parse(`{{template "authboss" .}}`))

	ab.Cfg.XSRFName = "csrf_token"
	ab.Cfg.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	ab.Cfg.CookieStoreMaker = NewCookieStorer
	ab.Cfg.SessionStoreMaker = NewSessionStorer

	//ab.Cfg.Mailer = ab.SMTPMailer("smtp.gmail.com:587", smtp.PlainAuth("you@gmail.com", "you@gmail.com", "password", "smtp.gmail.com"))
	ab.Cfg.Mailer = ab.LogMailer(os.Stdout)

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

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

	mux.PathPrefix("/auth").Handler(ab.NewRouter())

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
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "Not found")
	})

	// Set up our middleware chain
	stack := alice.New(logger, nosurfing, ab.ExpireMiddleware).Then(mux)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Println(http.ListenAndServe("localhost:"+port, stack))
}

func layoutData(w http.ResponseWriter, r *http.Request) ab.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return ab.HTMLData{
		"loggedin":          userInter != nil,
		"username":          "",
		ab.FlashSuccessKey:  ab.FlashSuccess(w, r),
		ab.FlashErrorKey:    ab.FlashError(w, r),
		"csrf_token":        nosurf.Token(r),
		"current_user_name": currentUserName,
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

	data := layoutData(w, r).MergeKV("post", blogs.Get(id))
	err := templates.Render(w, "show", data)
	if err != nil {
		panic(err)
	}
}

func new(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("post", Blog{})
	mustRender(w, "new", data)
}

var nextID = len(blogs) + 1

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

	b.ID = nextID
	nextID++
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

	data := layoutData(w, r).MergeKV("post", blogs.Get(id))
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

	var b = blogs.Get(id)
	if badRequest(w, schemaDec.Decode(b, r.PostForm)) {
		return
	}

	b.Date = time.Now()

	data := layoutData(w, r).MergeKV("post", b)
	mustRender(w, "show", data)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

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

	if id <= 0 {
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
