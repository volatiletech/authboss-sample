package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/justinas/nosurf"
)

type authProtector struct {
	f http.HandlerFunc
}

func authProtect(f func(w http.ResponseWriter, r *http.Request)) authProtector {
	return authProtector{http.HandlerFunc(f)}
}

func (ap authProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := ab.CurrentUser(w, r); err != nil {
		log.Println("Error fetching current user:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		log.Println("Redirecting unauthorized user from:", r.URL.Path)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		ap.f(w, r)
	}
}

func nosurfing(h http.Handler) http.Handler {
	surfing := nosurf.New(h)
	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Failed to validate XSRF Token:", nosurf.Reason(r))
		w.WriteHeader(http.StatusBadRequest)
	}))
	return surfing
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		session, err := sessionStore.Get(r, sessionCookieName)
		if err == nil {
			fmt.Print("Session: ")
			first := true
			for k, v := range session.Values {
				if first {
					first = false
				} else {
					fmt.Print(", ")
				}
				fmt.Printf("%s = %v", k, v)
			}
			fmt.Println()
		}
		fmt.Print("Database: ")
		for _, u := range database.Users {
			fmt.Printf("%#v\n", u)
		}
		fmt.Printf("Context: %s\n", spew.Sdump(r.Context()))
		h.ServeHTTP(w, r)
	})
}
