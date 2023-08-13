package model

type InstanceRequest struct {
  Id string
  Type string
  Src string
}

type InstanceResponse struct {
  Id string
  Status string
  Stdout string
  Stderr string
}
