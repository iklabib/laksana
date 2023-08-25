package main

type Box struct {
	Id      string
	Type    string
	Version string
}

// TODO: show compile error
type Response struct {
	Status string
	Stdout string
	Stderr string
}
