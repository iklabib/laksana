package model

type Box struct {
  Id string
  Type string
  Version string
}

// TODO: show compile error
type Response struct {
  Id string
  Status string
  Stdout string
  Stderr string
}
