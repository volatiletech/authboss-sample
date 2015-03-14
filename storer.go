package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/authboss.v0"
)

var nextUserID int

type User struct {
	ID   int
	Name string

	// Auth
	Email    string
	Password string

	// OAuth2
	Oauth2Uid      string
	Oauth2Provider string
	Oauth2Token    string
	Oauth2Refresh  string
	Oauth2Expiry   time.Time

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Lock
	AttemptNumber int
	AttemptTime   time.Time
	Locked        bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time

	// Remember is in another table
}

type MemStorer struct {
	Users  map[string]User
	Tokens map[string]string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"zeratul@heroes.com": User{
				ID:        1,
				Name:      "Zeratul",
				Password:  "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:     "zeratul@heroes.com",
				Confirmed: true,
			},
		},
		Tokens: make(map[string]string),
	}
}

func (s MemStorer) Create(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	user.ID = nextUserID
	nextUserID++

	s.Users[key] = user
	fmt.Println("Create")
	spew.Dump(s.Users)
	return nil
}

func (s MemStorer) Put(key string, attr authboss.Attributes) error {
	return s.Create(key, attr)
}

func (s MemStorer) Get(key string, attrMeta authboss.AttributeMeta) (result interface{}, err error) {
	user, ok := s.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s MemStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s MemStorer) GetOAuth(uid, provider string, attrMeta authboss.AttributeMeta) (result interface{}, err error) {
	user, ok := s.Users[uid+provider]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s MemStorer) AddToken(key, token string) error {
	s.Tokens[key] = token
	fmt.Println("AddToken")
	spew.Dump(s.Tokens)
	return nil
}

func (s MemStorer) DelTokens(key string) error {
	delete(s.Tokens, key)
	fmt.Println("DelTokens")
	spew.Dump(s.Tokens)
	return nil
}

func (s MemStorer) UseToken(givenKey, token string) (key string, err error) {
	t, ok := s.Tokens[givenKey]
	if !ok {
		return "", authboss.ErrTokenNotFound
	}

	s.DelTokens(givenKey)
	return t, nil
}

func (s MemStorer) ConfirmUser(tok string) (result interface{}, err error) {
	fmt.Println("==============", tok)

	for _, u := range s.Users {
		if u.ConfirmToken == tok {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}

func (s MemStorer) RecoverUser(rec string) (result interface{}, err error) {
	for _, u := range s.Users {
		if u.RecoverToken == rec {
			return &u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}
