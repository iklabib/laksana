// TODO: should be it's own buildable
package instance

import (
	"encoding/json"
	model "markisa/model"
	"net/http"
	"os"
	"os/exec"
)

func task(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  if (r.Method == "POST") {
    var req model.InstanceRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    }

    res := model.InstanceResponse {
      Id: req.Id,
    }

    f, err := os.Create("main.py")
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    } 
    defer f.Close()

    f.Write([]byte(req.Src))

    cmd := exec.Command("python3", "./main.py")
    if out, err := cmd.Output(); err != nil {
      res.Stderr = err.Error()
      res.Status = "ERROR"
    } else {
      res.Stdout = string(out)
      res.Status = "DONE"
    }

    result, err := json.Marshal(res)
    if err != nil {
      http.Error(w, err.Error(), 500)
      return
    } 
    if _, err := w.Write(result); err != nil {
      panic(err)
    }
    return 
  }
}

func main() {
  http.HandleFunc("/", task)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}
