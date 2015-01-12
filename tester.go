package main

import (
	"net/http"
	"os"

	"log"

	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"html/template"
)

type MongoStorer struct {
	users *mgo.Collection
}

type MGOUser struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

func (s MongoStorer) Create(key string, attr authboss.Attributes) error {
	return nil
}

func (s MongoStorer) Put(key string, attr authboss.Attributes) error {
	return nil
}

func (s MongoStorer) Get(key string, attrMeta authboss.AttributeMeta) (result interface{}, err error) {
	u := MGOUser{}
	err = s.users.Find(bson.M{"username": key}).One(&u)
	return u, err
}

func main() {
	c := authboss.NewConfig()

	if session, err := mgo.Dial("authboss:authboss@localhost/authboss"); err != nil {
		log.Fatal(err)
	} else {
		mgo := session.DB("authboss")
		c.Storer = &MongoStorer{mgo.C("users")}
	}

	c.LogWriter = os.Stdout
	c.ViewsPath = "views"
	c.AuthLoginSuccessRoute = "/dashboard"

	if err := authboss.Init(c); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", authboss.NewRouter(c))

	templates, _ := template.ParseFiles("views/dashboard.tpl")
	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, _ *http.Request) {
		templates.ExecuteTemplate(w, "dashboard.tpl", nil)
	})

	http.ListenAndServe("localhost:8080", mux)
}
