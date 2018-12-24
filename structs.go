package main

type Idol struct {
	Rank       int
	Name       string
	VoteAmount int
}

type Role struct {
	Name  string
	Idols []Idol
}

type Theme struct {
	Name  string
	Roles []Role
}
