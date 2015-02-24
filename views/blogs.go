package main

import "time"

type Blog struct {
	Title    string
	AuthorID string
	Date     time.Time
}

type Blogs []Blog
