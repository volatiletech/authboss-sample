package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

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
		b, err := httputil.DumpRequest(r, true)
		if err != nil {
			fmt.Println("What:", err)
		}
		os.Stdout.Write(b)
		h.ServeHTTP(w, r)
	})
}

func touch(h http.Handler) http.Handler {
	return expire.Middleware(NewSessionStorer, h)
}
