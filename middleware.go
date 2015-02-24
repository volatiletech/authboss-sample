package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

func nosurfHandler(h http.Handler) http.Handler {
	return nosurf.New(h)
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := httputil.DumpRequest(r)
		if err != nil {
			fmt.Println("What:", err)
		}
		os.Stdout.Write(b)
	})
}
