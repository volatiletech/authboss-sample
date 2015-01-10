package main

import (
	"net/http"
	"os"

	"log"

	"gopkg.in/authboss.v0"
	_ "gopkg.in/authboss.v0/auth"
)

type MongoStorer struct {
}

func (s MongoStorer) Create(key string, attr authboss.Attributes) error {
	return nil
}

func (s MongoStorer) Put(key string, attr authboss.Attributes) error {
	return nil
}

func (s MongoStorer) Get(key string, attrMeta authboss.AttributeMeta) (interface{}, error) {
	return nil, nil
}

func main() {
	c := authboss.NewConfig()

	c.Storer = &MongoStorer{}
	c.LogWriter = os.Stdout

	if err := authboss.Init(c); err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe("localhost:8080", authboss.Router(c))
}
