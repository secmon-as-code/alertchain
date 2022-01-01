package model

type Reference struct {
	Source  string
	Title   string
	URL     string
	Comment string
}

type References []*Reference
