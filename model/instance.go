package model

type Box struct {
	Id      string
	Type    string
	Version string
}

type Request struct {
  Token string
  Archive string
  Type string
}

// TODO: show compile error
type Response struct {
	Status string
	Stdout string
	Stderr string
}
