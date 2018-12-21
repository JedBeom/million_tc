package main

type Idol struct {
	Rank       int
	Name       string
	VoteAmount int
}

type Role struct {
	Name  string
	idols []Idol
}

type Theme struct {
	Name  string
	roles []Role
}
