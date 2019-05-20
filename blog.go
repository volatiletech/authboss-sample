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
	_ "github.com/volatiletech/authboss/auth"
	"github.com/volatiletech/authboss/confirm"
	"github.com/volatiletech/authboss/defaults"
	"github.com/volatiletech/authboss/lock"
	_ "github.com/volatiletech/authboss/logout"
	aboauth "github.com/volatiletech/authboss/oauth2"
	"github.com/volatiletech/authboss/otp/twofactor"
	"github.com/volatiletech/authboss/otp/twofactor/sms2fa"
	"github.com/volatiletech/authboss/otp/twofactor/totp2fa"
	_ "github.com/volatiletech/authboss/recover"
	_ "github.com/volatiletech/authboss/register"
	"github.com/volatiletech/authboss/remember"

	"github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss-renderer"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/aarondl/tpl"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
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
	ab        = authboss.New()
	database  = NewMemStorer()
	schemaDec = schema.NewDecoder()

	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer

	templates tpl.Templates
)

const (
	sessionCookieName = "ab_blog"
)

func setupAuthboss() {
	ab.Config.Paths.RootURL = "http://localhost:3000"

	if !*flagAPI {
		// Prevent us from having to use Javascript in our basic HTML
		// to create a delete method, but don't override this default for the API
		// version
		ab.Config.Modules.LogoutMethod = "GET"
	}

	// Set up our server, session and cookie storage mechanisms.
	// These are all from this package since the burden is on the
	// implementer for these.
	ab.Config.Storage.Server = database
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore

	// Another piece that we're responsible for: Rendering views.
	// Though note that we're using the authboss-renderer package
	// that makes the normal thing a bit easier.
	if *flagAPI {
		ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	} else {
		ab.Config.Core.ViewRenderer = abrenderer.NewHTML("/auth", "ab_views")
	}

	// We render mail with the authboss-renderer but we use a LogMailer
	// which simply sends the e-mail to stdout.
	ab.Config.Core.MailRenderer = abrenderer.NewEmail("/auth", "ab_views")

	// The preserve fields are things we don't want to
	// lose when we're doing user registration (prevents having
	// to type them again)
	ab.Config.Modules.RegisterPreserveFields = []string{"email", "name"}

	// TOTP2FAIssuer is the name of the issuer we use for totp 2fa
	ab.Config.Modules.TOTP2FAIssuer = "ABBlog"
	ab.Config.Modules.RoutesRedirectOnUnauthed = true

	// Turn on e-mail authentication required
	ab.Config.Modules.TwoFactorEmailAuthRequired = true

	// This instantiates and uses every default implementation
	// in the Config.Core area that exist in the defaults package.
	// Just a convenient helper if you don't want to do anything fancy.
	defaults.SetCore(&ab.Config, *flagAPI, false)

	// Here we initialize the bodyreader as something customized in order to accept a name
	// parameter for our user as well as the standard e-mail and password.
	//
	// We also change the validation for these fields
	// to be something less secure so that we can use test data easier.
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]{1,}`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: *flagAPI,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
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

	// Set up 2fa
	twofaRecovery := &twofactor.Recovery{Authboss: ab}
	if err := twofaRecovery.Setup(); err != nil {
		panic(err)
	}

	totp := &totp2fa.TOTP{Authboss: ab}
	if err := totp.Setup(); err != nil {
		panic(err)
	}

	sms := &sms2fa.SMS{Authboss: ab, Sender: smsLogSender{}}
	if err := sms.Setup(); err != nil {
		panic(err)
	}

	// Set up Google OAuth2 if we have credentials in the
	// file oauth2.toml for it.
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

	// Initialize authboss (instantiate modules etc.)
	if err := ab.Init(); err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	// Load our application's templates
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
	cookieStore = abclientstate.NewCookieStorer(cookieStoreKey, nil)
	cookieStore.HTTPOnly = false
	cookieStore.Secure = false
	sessionStore = abclientstate.NewSessionStorer(sessionCookieName, sessionStoreKey, nil)
	cstore := sessionStore.Store.(*sessions.CookieStore)
	cstore.Options.HttpOnly = false
	cstore.Options.Secure = false
	cstore.MaxAge(int((30 * 24 * time.Hour) / time.Second))

	// Initialize authboss
	setupAuthboss()

	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)

	mux := chi.NewRouter()
	// The middlewares we're using:
	// - logger just does basic logging of requests and debug info
	// - nosurfing is a more verbose wrapper around csrf handling
	// - LoadClientStateMiddleware is required for session/cookie stuff
	// - remember middleware logs users in if they have a remember token
	// - dataInjector is for putting data into the request context we need for our template layout
	mux.Use(logger, nosurfing, ab.LoadClientStateMiddleware, remember.Middleware(ab), dataInjector)

	// Authed routes
	mux.Group(func(mux chi.Router) {
		mux.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized), lock.Middleware(ab), confirm.Middleware(ab))
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

	if *flagAPI {
		// In order to have a "proper" API with csrf protection we allow
		// the options request to return the csrf token that's required to complete the request
		// when using post
		optionsHandler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-CSRF-TOKEN", nosurf.Token(r))
			w.WriteHeader(http.StatusOK)
		}

		// We have to add each of the authboss get/post routes specifically because
		// chi sees the 'Mount' above as overriding the '/*' pattern.
		routes := []string{"login", "logout", "recover", "recover/end", "register"}
		mux.MethodFunc("OPTIONS", "/*", optionsHandler)
		for _, r := range routes {
			mux.MethodFunc("OPTIONS", "/auth/"+r, optionsHandler)
		}
	}

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
		r.Body.Close()
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
		r.Body.Close()
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

type smsLogSender struct {
}

// Send an SMS
func (s smsLogSender) Send(ctx context.Context, number, text string) error {
	fmt.Println("sms sent to:", number, "contents:", text)
	return nil
}
