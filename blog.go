package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/volatiletech/authboss"
	// Auth module for authboss
	"github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/defaults"

	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var (
	ab        = authboss.New()
	database  = NewMemStorer()
	templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	schemaDec = schema.NewDecoder()
)

func setupAuthboss() {
	ab.Config.Modules.AuthLogoutMethod = "GET"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = NewSessionStorer()
	ab.Config.Storage.CookieState = NewCookieStorer()

	ab.Config.Core.ViewRenderer = abrenderer.New("/auth")

	defaults.SetDefaultCore(&ab.Config, false)

	if err := ab.Init(); err != nil {
		panic(err)
	}

	/*ab.Storer = database
	ab.OAuth2Storer = database
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views"
	ab.RootURL = `http://localhost:3000`

	ab.LayoutDataMaker = layoutData

	ab.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"google": authboss.OAuth2Provider{
			OAuth2Config: &oauth2.Config{
				ClientID:     ``,
				ClientSecret: ``,
				Scopes:       []string{`profile`, `email`},
				Endpoint:     google.Endpoint,
			},
			Callback: aboauth.Google,
		},
	}

	b, err := ioutil.ReadFile(filepath.Join("views", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b)))

	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	ab.CookieStoreMaker = NewCookieStorer
	ab.SessionStoreMaker = NewSessionStorer

	ab.Mailer = authboss.LogMailer(os.Stdout)

	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       8,
			AllowWhitespace: false,
		},
	}

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}*/
}

func main() {
	// Initialize Sessions and Cookies
	// Typically gorilla securecookie and sessions packages require
	// highly random secret keys that are not divulged to the public.
	//
	// In this example we use keys generated one time (if these keys ever become
	// compromised the gorilla libraries allow for key rotation, see gorilla docs)
	// The keys are 64-bytes as recommended for HMAC keys as per the gorilla docs.
	//
	// These values MUST be changed for any new project as these keys are already "compromised"
	// as they're in the public domain, if you do not change these your application will have a fairly
	// wide-opened security hole. You can generate your own with the code below, or using whatever method
	// you prefer:
	//
	//    func main() {
	//        fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)))
	//    }
	//
	// We store them in base64 in the example to make it easy if we wanted to move them later to
	// a configuration environment var or file.
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`NpEPi8pEjKVjLGJ6kYCS+VTCzi6BUuDzU0wrwXyf5uDPArtlofn2AG6aTMiPmN3C909rsEWMNqJqhIVPGP3Exg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`AbfYwmmt8UCwUuhd9qvfNA9UCuN1cVcKJN1ofbiky6xCyyBj20whe40rJa3Su0WOWLWcPpO1taqJdsEI/65+JA==`)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)

	// Initialize ab.
	setupAuthboss()

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)

	mux := chi.NewRouter()
	mux.Use(ab.ClientStateMiddleware, logger) //, nosurfing)

	// Routes
	mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))

	mux.Method("GET", "/blogs/new", authProtect(newblog))
	mux.Method("GET", "/blogs/{id}/edit", authProtect(edit))

	mux.Get("/blogs", index)
	mux.Get("/", index)

	mux.Method("POST", "/blogs/{id}/edit", authProtect(update))
	mux.Method("POST", "/blogs/new", authProtect(create))

	// This should actually be a DELETE but I can't be bothered to make a proper
	// destroy link using javascript atm.
	mux.Method("GET", "/blogs/{id}/destroy", authProtect(destroy))

	// Set up our middleware chain

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	log.Println(http.ListenAndServe("localhost:"+port, mux))
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":          userInter != nil,
		"username":          "",
		"current_user_name": currentUserName,
	}
	/*
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}*/
}

func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, r, "index", data)
}

func newblog(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("post", Blog{})
	mustRender(w, r, "new", data)
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
	b.AuthorID = "Rick"

	blogs = append(blogs, b)

	http.Redirect(w, r, "/", http.StatusFound)
}

func edit(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	data := layoutData(w, r).MergeKV("post", blogs.Get(id))
	mustRender(w, r, "edit", data)
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

	http.Redirect(w, r, "/", http.StatusFound)
}

func destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

	http.Redirect(w, r, "/", http.StatusFound)
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

func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	//data.MergeKV("csrf_token", nosurf.Token(r))
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
