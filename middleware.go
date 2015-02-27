package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
	"gopkg.in/authboss.v0/expire"
)

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
		session, err := sessionStore.Get(r, "derpasaurous")
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
			fmt.Printf("%s: Confirmed: %v ConfirmToken: %v AttemptN: %v AttemptT: %v Locked: %v RecoverTok: %v RecoverExp: %v\n",
				u.Email, u.Confirmed, u.ConfirmToken, u.AttemptNumber, u.AttemptTime, u.Locked, u.RecoverToken, u.RecoverTokenExpiry)
		}
		h.ServeHTTP(w, r)
	})
}

func touch(h http.Handler) http.Handler {
	return expire.Middleware(NewSessionStorer, h)
}
