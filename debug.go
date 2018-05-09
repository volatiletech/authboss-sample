package main

import "fmt"

func debugln(args ...interface{}) {
	if !*flagDebug {
		return
	}

	fmt.Println(args...)
}

func debugf(format string, args ...interface{}) {
	if !*flagDebug {
		return
	}

	fmt.Printf(format, args...)
}
