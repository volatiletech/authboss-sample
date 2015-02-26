package main

import "time"

type Blog struct {
	ID       int
	Title    string
	AuthorID string
	Date     time.Time
	Content  string
}

type Blogs []Blog

var blogs = Blogs{
	{1, "My first Hook", "Stitches", time.Now().AddDate(0, 0, -3),
		`When I was young I was weak. But then I grew up and now I'm ridiculous` +
			`because I can hook a whole team and destroy them with gorge lololo. ` +
			`Look at me mom. I'm a big kid now!`,
	},
	{2, "Halp the nerfed!11", "Murky", time.Now().AddDate(0, 0, -1),
		`I used to be really amazing, then they nerfed me, then I was really` +
			`good in the right hands, now... I have no idea, why Blizzard? Why??!!`,
	},
}

func (blgs *Blogs) Get(id int) *Blog {
	for i, _ := range blogs {
		b := &blogs[i]
		if b.ID == id {
			return b
		}
	}
	return nil
}

func (blgs *Blogs) Delete(id int) {
	if len(blogs) == 1 {
		blogs = []Blog{}
		return
	}

	found := -1
	for i, _ := range blogs {
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
