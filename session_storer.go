package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/volatiletech/authboss"
)

const (
	sessionCookieName = "ab_blog"
)

var (
	sessionStore *sessions.CookieStore
)

// SessionState for sessions
type SessionState struct {
	session *sessions.Session
}

// Get a key from the session
func (s SessionState) Get(key string) (string, bool) {
	str, ok := s.session.Values[key]
	if !ok {
		return "", false
	}
	value := str.(string)
	debugf("Got session (%s): %s\n", key, value)

	return value, ok
}

// SessionStorer stores sessions in a global gorilla cookiestore
type SessionStorer struct{}

// NewSessionStorer constructor
func NewSessionStorer() *SessionStorer {
	return &SessionStorer{}
}

// ReadState loads the session from the request context
func (s SessionStorer) ReadState(r *http.Request) (authboss.ClientState, error) {
	debugln("Loading session state")
	session, err := sessionStore.Get(r, sessionCookieName)
	if err != nil {
		return nil, nil
	}

	if session == nil {
		debugln("Session was nil")
	}

	cs := &SessionState{
		session: session,
	}

	return cs, nil
}

// WriteState to the responsewriter
func (s SessionStorer) WriteState(w http.ResponseWriter, state authboss.ClientState, ev []authboss.ClientStateEvent) error {
	debugln("Writing session state")
	ses := state.(*SessionState)

	for _, ev := range ev {
		switch ev.Kind {
		case authboss.ClientStateEventPut:
			ses.session.Values[ev.Key] = ev.Value
		case authboss.ClientStateEventDel:
			delete(ses.session.Values, ev.Key)
		}
	}

	return sessionStore.Save(nil, w, ses.session)
}

/*
func (s SessionStorer) Get(key string) (string, bool) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
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
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values[key] = value
	session.Save(s.r, s.w)
}

func (s SessionStorer) Del(key string) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(session.Values, key)
	session.Save(s.r, s.w)
}
*/
