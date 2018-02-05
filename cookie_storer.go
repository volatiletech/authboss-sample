package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/volatiletech/authboss"
)

var cookieStore *securecookie.SecureCookie

// Cookies is a struct to hold cookies for the duration of the request
type Cookies struct {
	cookies map[string]*http.Cookie
}

// Get a cookie's value
func (c Cookies) Get(key string) (string, bool) {
	if cookie, ok := c.cookies[key]; ok {
		var value string
		if err := cookieStore.Decode(cookie.Name, cookie.Value, &value); err != nil {
			panic("COOKIE DECODE FAILURE:" + err.Error())
		}

		return value, true
	}

	return "", false
}

// CookieStorer writes and reads cookies
type CookieStorer struct{}

// NewCookieStorer constructor
func NewCookieStorer() *CookieStorer {
	return &CookieStorer{}
}

// ReadState from the request
func (c CookieStorer) ReadState(w http.ResponseWriter, r *http.Request) (authboss.ClientState, error) {
	cs := &Cookies{
		cookies: make(map[string]*http.Cookie),
	}

	for _, c := range r.Cookies() {
		cs.cookies[c.Name] = c
	}

	return cs, nil
}

// WriteState to the responsewriter
func (c CookieStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	for _, ev := range ev {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			encoded, err := cookieStore.Encode(ev.Key, ev.Value)
			if err != nil {
				fmt.Println(err)
				return err
			}

			cookie := &http.Cookie{
				Expires: time.Now().UTC().AddDate(1, 0, 0),
				Name:    ev.Key,
				Value:   encoded,
				Path:    "/",
			}
			http.SetCookie(w, cookie)
		case authboss.ClientStateEventDel:
			cookie := &http.Cookie{
				MaxAge: -1,
				Name:   ev.Key,
				Path:   "/",
			}
			http.SetCookie(w, cookie)
		}
	}

	return nil
}
