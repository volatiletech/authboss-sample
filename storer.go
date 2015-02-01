package main

import (
	"time"

	"gopkg.in/authboss.v0"
)

type User struct {
	Username           string    `bson:"username"`
	Password           string    `bson:"password"`
	Email              string    `bson:"email"`
	RecoverToken       string    `bson:"recoverToken"`
	RecoverTokenExpiry time.Time `bson:"recoverTokenExpiry"`
}

type MemStorer struct {
	Users  map[string]User
	Tokens map[string]string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"kris": User{"kris", "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", "kris@test.com", "", time.Now().UTC()}, // pass = 1234
		},
		Tokens: make(map[string]string),
	}
}

func (s MemStorer) Create(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user); err != nil {
		return err
	}

	s.Users[key] = user
	//spew.Dump(s.Users)
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

	return user, nil
}

func (s MemStorer) AddToken(key, token string) error {
	s.Tokens[key] = token
	//spew.Dump(s.Tokens)
	return nil
}

func (s MemStorer) DelTokens(key string) error {
	delete(s.Tokens, key)
	//spew.Dump(s.Tokens)
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

func (s MemStorer) RecoverUser(rec string) (result interface{}, err error) {
	for _, u := range s.Users {
		if u.RecoverToken == rec {
			return u, nil
		}
	}

	return nil, authboss.ErrUserNotFound
}
