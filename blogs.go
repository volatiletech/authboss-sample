package main

import "time"

// Blog data
type Blog struct {
	ID       int       `json:"id,omitempty"`
	Title    string    `json:"title,omitempty"`
	AuthorID string    `json:"author_id,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Content  string    `json:"content,omitempty"`
}

// Blogs storage
type Blogs []Blog

var blogs = Blogs{
	{1, "My first portal", "Rick", time.Now().AddDate(0, 0, -3).Add(-time.Hour * 2),
		`I successfully opened a portal to another dimension, I think it's pretty clear that I'm the smartest person on earth ` +
			`and this'll let me go see if there's anything out in the verse that can compete with my tremendous intellect, after ` +
			`dragging Morty along on a few adventures I think the answer is still a resounding: no.`,
	},
	{2, "My Life", "Morty", time.Now().AddDate(0, 0, -1),
		`My Grandpa is a really cool guy, but who I really think is great is Jessica. I keep staring at her in class hoping ` +
			`that one day she'll realize just how great a guy she's missing out on. She doesn't need any of these bad ` +
			`guys she keeps dating, that's only going to hurt her. I'm a whole lot of Morty, and I'm waiting for you Jessica.`,
	},
}

// Get blogs
func (blgs *Blogs) Get(id int) *Blog {
	for i := range blogs {
		b := &blogs[i]
		if b.ID == id {
			return b
		}
	}
	return nil
}

// Delete a blog
func (blgs *Blogs) Delete(id int) {
	if len(blogs) == 1 {
		blogs = []Blog{}
		return
	}

	found := -1
	for i := range blogs {
		b := &blogs[i]
		if b.ID == id {
			found = i
		}
	}

	for i := found; i < len(blogs)-1; i++ {
		blogs[i], blogs[i+1] = blogs[i+1], blogs[i]
	}
	blogs = blogs[:len(blogs)-1]
}
