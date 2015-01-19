package main // import "gopkg.in/authboss-sample.v0"

import (
	"net/http"
	"os"

	"log"

	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	_ "gopkg.in/authboss.v0/recover"
	_ "gopkg.in/authboss.v0/remember"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"html/template"

	"fmt"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var cookieStore *securecookie.SecureCookie
var sessionStore *sessions.CookieStore

type CookieStorer struct {
	w http.ResponseWriter
	r *http.Request
}

func NewCookieStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &CookieStorer{w, r}
}

func (s CookieStorer) Get(key string) (string, bool) {
	cookie, err := s.r.Cookie(key)
	if err != nil {
		return "", false
	}

	var value string
	err = cookieStore.Decode(key, cookie.Value, &value)
	if err != nil {
		return "", false
	}

	return value, true
}

func (s CookieStorer) Put(key, value string) {
	encoded, err := cookieStore.Encode(key, value)
	if err != nil {
		fmt.Println(err)
	}

	cookie := &http.Cookie{
		Name:  key,
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(s.w, cookie)
}

func (s CookieStorer) Del(key string) {
}

type SessionStorer struct {
	w http.ResponseWriter
	r *http.Request
}

func NewSessionStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &SessionStorer{w, r}
}

func (s SessionStorer) Get(key string) (string, bool) {
	session, err := sessionStore.Get(s.r, "derpasaurous")
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	strInf, ok := session.Values[key]
	if !ok {
		return "", false
	}

	str, ok := strInf.(string)
	if !ok {
		return "", false
	}

	return str, true
}

func (s SessionStorer) Put(key, value string) {
	session, err := sessionStore.Get(s.r, "derpasaurous")
	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values[key] = value
	session.Save(s.r, s.w)
}

func (s SessionStorer) Del(key string) {
	session, err := sessionStore.Get(s.r, "derpasaurous")
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(session.Values, key)
	session.Save(s.r, s.w)
}

type MongoStorer struct {
	users  *mgo.Collection
	tokens *mgo.Collection
}

type MGOUser struct {
	Username   string `bson:"username"`
	Password   string `bson:"password"`
	Email      string `bson:"email"`
	Resettoken string `bson:"resetToken"`
}

type MGOToken struct {
	Username string `bson:"username"`
	Token    string `bson:"token"`
}

func (s MongoStorer) Create(key string, attr authboss.Attributes) error {
	return nil
}

func (s MongoStorer) Put(key string, attr authboss.Attributes) error {
	u := MGOUser{}
	if err := s.users.Find(bson.M{"username": key}).One(&u); err != nil {
		return err
	}

	if err := attr.Bind(&u); err != nil {
		log.Println("error", err)
	}

	if err := s.users.Update(bson.M{"username": key}, u); err != nil {
		return err
	}

	return nil
}

func (s MongoStorer) Get(key string, attrMeta authboss.AttributeMeta) (result interface{}, err error) {
	u := MGOUser{}
	err = s.users.Find(bson.M{"username": key}).One(&u)
	return u, err
}

func (s MongoStorer) AddToken(key, token string) error {
	t := MGOToken{key, token}
	return s.tokens.Insert(t)
}

func (s MongoStorer) DelTokens(key string) error {
	return nil
}

func (s MongoStorer) UseToken(givenKey, token string) (key string, err error) {
	t := MGOToken{}
	sel := bson.M{"username": givenKey, "token": token}

	if err = s.tokens.Find(sel).One(&t); err != nil {
		return "", authboss.TokenNotFound
	}

	if err = s.tokens.Remove(sel); err != nil {
		return "", err
	}

	return t.Username, nil
}

func main() {
	c := authboss.NewConfig()
	cookieStore = securecookie.New([]byte("very-secret"), nil)
	sessionStore = sessions.NewCookieStore([]byte("asdf"))

	if session, err := mgo.Dial("authboss:authboss@localhost/authboss"); err != nil {
		log.Fatal(err)
	} else {
		mgo := session.DB("authboss")
		c.Storer = &MongoStorer{mgo.C("users"), mgo.C("tokens")}
	}

	c.LogWriter = os.Stdout
	c.ViewsPath = "views"
	c.AuthLoginSuccessRoute = "/dashboard"
	c.CookieStoreMaker = NewCookieStorer
	c.SessionStoreMaker = NewSessionStorer

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
