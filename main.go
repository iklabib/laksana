package main

import (
  "fmt"
  task "markisa/task"
  util "markisa/util"
	container "markisa/container"
	"markisa/model"
)


func main() {
  port, err := util.GetFreePort()
  if err != nil {
    panic(err)
  }

  var hostPort uint16 = uint16(port)
  var containerPort uint16 = 8080

  client := container.Init()
  id := client.Create("localhost/kardus:python", hostPort, containerPort)
  client.Run(id)
  
  req := model.InstanceRequest{
    Id: id,
    Type: "python",
    Src: `for it in range(5):print(f"counting down: {it+1}")`,
  }

  url := fmt.Sprintf("http://localhost:%d", port)
  task := task.SendTask(&req, url)

  fmt.Println("ID: " + task.Id)
  fmt.Println("Status: " + task.Status)
  fmt.Println("Stdout: " + task.Stdout)
  fmt.Println("Stderr: " + task.Stderr)

  client.Stop(id)
  client.Remove(id)
}
