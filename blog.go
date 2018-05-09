package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss-renderer"
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/confirm"
	"github.com/volatiletech/authboss/defaults"
	"github.com/volatiletech/authboss/lock"
	_ "github.com/volatiletech/authboss/logout"
	aboauth "github.com/volatiletech/authboss/oauth2"
	_ "github.com/volatiletech/authboss/recover"
	_ "github.com/volatiletech/authboss/register"
	"github.com/volatiletech/authboss/remember"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/aarondl/tpl"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var (
	flagDebug    = flag.Bool("debug", false, "output debugging information")
	flagDebugDB  = flag.Bool("debugdb", false, "output database on each request")
	flagDebugCTX = flag.Bool("debugctx", false, "output specific authboss related context keys on each request")
	flagAPI      = flag.Bool("api", false, "configure the app to be an api instead of an html app")
)

var (
	ab                      = authboss.New()
	database                = NewMemStorer()
	schemaDec               = schema.NewDecoder()
	templates tpl.Templates = nil
)

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:3000"
	ab.Config.Modules.LogoutMethod = "GET"

	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = NewSessionStorer()
	ab.Config.Storage.CookieState = NewCookieStorer()

	ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")
	ab.Config.Core.Mailer = defaults.LogMailer{}

	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	defaults.SetCore(&ab.Config, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}

	ab.Config.Core.BodyReader = defaults.HTTPFormReader{
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule},
			"recover_end": {passwordRule},
		},
		Confirms: map[string][]string{
			"register":    {"password", authboss.ConfirmPrefix + "password"},
			"recover_end": {"password", authboss.ConfirmPrefix + "password"},
		},
		Whitelist: map[string][]string{
			"register": []string{"email", "name", "password"},
		},
	}

	oauthcreds := struct {
		ClientID     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	}{}

	_, err := toml.DecodeFile("oauth2.toml", &oauthcreds)
	if err == nil && len(oauthcreds.ClientID) != 0 && len(oauthcreds.ClientSecret) != 0 {
		fmt.Println("oauth2.toml exists, configuring google oauth2")
		ab.Config.Modules.OAuth2Providers = map[string]authboss.OAuth2Provider{
			"google": authboss.OAuth2Provider{
				OAuth2Config: &oauth2.Config{
					ClientID:     oauthcreds.ClientID,
					ClientSecret: oauthcreds.ClientSecret,
					Scopes:       []string{`profile`, `email`},
					Endpoint:     google.Endpoint,
				},
				FindUserDetails: aboauth.GoogleUserDetails,
			},
		}
	} else if os.IsNotExist(err) {
		fmt.Println("oauth2.toml doesn't exist, not registering oauth2 handling")
	} else {
		fmt.Println("error loading oauth2.toml:", err)
	}

	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if !*flagAPI {
		templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	}

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
	mux.Use(logger, nosurfing, ab.LoadClientStateMiddleware, remember.Middleware(ab), dataInjector)

	// Authed routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware(ab), lock.Middleware(ab), confirm.Middleware(ab))
		mux.MethodFunc("GET", "/blogs/new", newblog)
		mux.MethodFunc("GET", "/blogs/{id}/edit", edit)
		mux.MethodFunc("POST", "/blogs/{id}/edit", update)
		mux.MethodFunc("POST", "/blogs/new", create)
		// This should actually be a DELETE but can't be bothered to make a proper
		// destroy link using javascript atm. See where AB allows you to configure
		// the logout HTTP method.
		mux.MethodFunc("GET", "/blogs/{id}/destroy", destroy)
	})

	// Routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.ModuleListMiddleware(ab))
		mux.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})
	mux.Get("/blogs", index)
	mux.Get("/", index)

	// Start the server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	log.Printf("Listening on localhost: %s", port)
	log.Println(http.ListenAndServe("localhost:"+port, mux))
}

func dataInjector(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := layoutData(w, &r)
		r = r.WithContext(context.WithValue(r.Context(), authboss.CTXKeyData, data))
		handler.ServeHTTP(w, r)
	})
}

// layoutData is passing pointers to pointers be able to edit the current pointer
// to the request. This is still safe as it still creates a new request and doesn't
// modify the old one, it just modifies what we're pointing to in our methods so
// we're able to skip returning an *http.Request everywhere
func layoutData(w http.ResponseWriter, r **http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.LoadCurrentUser(r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":          userInter != nil,
		"current_user_name": currentUserName,
		"csrf_token":        nosurf.Token(*r),
		"flash_success":     authboss.FlashSuccess(w, *r),
		"flash_error":       authboss.FlashError(w, *r),
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	mustRender(w, r, "index", authboss.HTMLData{"posts": blogs})
}

func newblog(w http.ResponseWriter, r *http.Request) {
	mustRender(w, r, "new", authboss.HTMLData{"post": Blog{}})
}

var nextID = len(blogs) + 1

func create(w http.ResponseWriter, r *http.Request) {
	// TODO: Validation

	var b Blog
	if *flagAPI {
		byt, err := ioutil.ReadAll(r.Body)
		if badRequest(w, err) {
			return
		}

		if badRequest(w, json.Unmarshal(byt, &b)) {
			return
		}
	} else {
		err := r.ParseForm()
		if badRequest(w, err) {
			return
		}

		if badRequest(w, schemaDec.Decode(&b, r.PostForm)) {
			return
		}
	}

	abuser := ab.CurrentUserP(r)
	user := abuser.(*User)

	b.ID = nextID
	nextID++
	b.Date = time.Now()
	b.AuthorID = user.Name

	blogs = append(blogs, b)

	if *flagAPI {
		w.WriteHeader(http.StatusOK)
		return
	}

	redirect(w, r, "/")
}

func edit(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	mustRender(w, r, "edit", authboss.HTMLData{"post": blogs.Get(id)})
}

func update(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	// TODO: Validation

	var b = blogs.Get(id)

	if *flagAPI {
		byt, err := ioutil.ReadAll(r.Body)
		if badRequest(w, err) {
			return
		}

		if badRequest(w, json.Unmarshal(byt, &b)) {
			return
		}
	} else {
		err := r.ParseForm()
		if badRequest(w, err) {
			return
		}
		if badRequest(w, schemaDec.Decode(b, r.PostForm)) {
			return
		}
	}

	b.Date = time.Now()

	if *flagAPI {
		w.WriteHeader(http.StatusOK)
		return
	}

	redirect(w, r, "/")
}

func destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

	if *flagAPI {
		w.WriteHeader(http.StatusOK)
		return
	}

	redirect(w, r, "/")
}

func blogID(w http.ResponseWriter, r *http.Request) (int, bool) {
	str := chi.RouteContext(r.Context()).URLParam("id")

	id, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error parsing blog id:", err)
		redirect(w, r, "/")
		return 0, false
	}

	if id <= 0 {
		redirect(w, r, "/")
		return 0, false
	}

	return id, true
}

func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	// We've sort of hijacked the authboss mechanism for providing layout data
	// for our own purposes. There's nothing really wrong with this but it looks magical
	// so here's a comment.
	var current authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)
	if dataIntf == nil {
		current = authboss.HTMLData{}
	} else {
		current = dataIntf.(authboss.HTMLData)
	}

	current.MergeKV("csrf_token", nosurf.Token(r))
	current.Merge(data)

	if *flagAPI {
		w.Header().Set("Content-Type", "application/json")

		byt, err := json.Marshal(current)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("failed to marshal json:", err)
			fmt.Fprintln(w, `{"error":"internal server error"}`)
		}

		w.Write(byt)
		return
	}

	err := templates.Render(w, name, current)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	if *flagAPI {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", path)
		w.WriteHeader(http.StatusFound)
		fmt.Fprintf(w, `{"path": %q}`, path)
		return
	}

	http.Redirect(w, r, path, http.StatusFound)
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	if *flagAPI {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"error":"bad request"}`, err)
		return true
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)
	return true
}
