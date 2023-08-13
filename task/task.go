package task

import (
	"bytes"
	"encoding/json"
  "time"
	"net/http"
	"markisa/model"
)

func SendTask(req *model.InstanceRequest, url string) model.InstanceResponse {
  reqJson, err := json.Marshal(req)
  if err != nil {
    panic(err)
  }

  buffer := bytes.NewBuffer(reqJson)
  client := http.Client {
    Timeout: 5 * time.Second,
  }

  resp, err := client.Post(url, "application/json", buffer)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()

  var respInstance model.InstanceResponse
  if err := json.NewDecoder(resp.Body).Decode(&respInstance); err != nil {
    panic(err)
  }

  return respInstance
}

