package main

import (
	"context"

	"github.com/volatiletech/authboss"
)

var nextUserID int

// User struct for authboss
type User struct {
	ID   int
	Name string

	// Auth
	Email    string
	Password string

	/*
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
		AttemptNumber int64
		AttemptTime   time.Time
		Locked        time.Time

		// Recover
		RecoverToken       string
		RecoverTokenExpiry time.Time

		// Remember is in another table
	*/
}

// PutPID into user
func (u *User) PutPID(ctx context.Context, pid string) error {
	u.Email = pid
	return nil
}

// PutPassword into user
func (u *User) PutPassword(ctx context.Context, password string) error {
	u.Password = password
	return nil
}

// GetPID from user
func (u User) GetPID(ctx context.Context) (string, error) {
	return u.Email, nil
}

// GetPassword from user
func (u User) GetPassword(ctx context.Context) (string, error) {
	return u.Password, nil
}

// MemStorer stores users in memory
type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

// NewMemStorer constructor
func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"rick@councilofricks.com": User{
				ID:       1,
				Name:     "Rick",
				Password: "$2a$10$XtW/BrS5HeYIuOCXYe8DFuInetDMdaarMUJEOg/VA/JAIDgw3l4aG", // pass = 1234
				Email:    "rick@councilofricks.com",
				//Confirmed: true,
			},
		},
		Tokens: make(map[string][]string),
	}
}

// ServerStorer represents the data store that's capable of loading users
// and giving them a context with which to store themselves.
type ServerStorer interface {
	// Load will look up the user based on the passed the PrimaryID
	Load(ctx context.Context, key string) (User, error)

	// Save persists the user in the database
	Save(ctx context.Context, user User) error
}

// Save the user
func (s MemStorer) Save(ctx context.Context, user authboss.User) error {
	u := user.(*User)
	s.Users[u.Email] = *u

	return nil
}

// Load the user
func (s MemStorer) Load(ctx context.Context, key string) (user authboss.User, err error) {
	u, ok := s.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &u, nil
}

/*
func (s MemStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s MemStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	user, ok := s.Users[uid+provider]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil
}

func (s MemStorer) AddToken(key, token string) error {
	s.Tokens[key] = append(s.Tokens[key], token)
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

func (s MemStorer) UseToken(givenKey, token string) error {
	toks, ok := s.Tokens[givenKey]
	if !ok {
		return authboss.ErrTokenNotFound
	}

	for i, tok := range toks {
		if tok == token {
			toks[i], toks[len(toks)-1] = toks[len(toks)-1], toks[i]
			s.Tokens[givenKey] = toks[:len(toks)-1]
			return nil
		}
	}

	return authboss.ErrTokenNotFound
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
*/
