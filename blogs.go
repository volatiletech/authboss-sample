package main

import "time"

type Blog struct {
	Title    string
	AuthorID string
	Date     time.Time
	Content  string
}

var blogs = []Blog{
	{"My first Hook", "Stitches", time.Now().AddDate(0, 0, -3),
		`When I was young I was weak. But then I grew up and now I'm ridiculous` +
			`because I can hook a whole team and destroy them with gorge lololo. ` +
			`Look at me mom. I'm a big kid now!`,
	},
	{"Halp the nerfed!11", "Murky", time.Now().AddDate(0, 0, -1),
		`I used to be really amazing, then they nerfed me, then I was really` +
			`good in the right hands, now... I have no idea, why Blizzard? Why??!!`,
	},
}
