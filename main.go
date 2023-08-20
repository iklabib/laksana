package main

import (
	"markisa/container"
	"fmt"
	"os"
)

func main() {
  buff, err := os.ReadFile("example.tar.gz")
  if err != nil {
    panic(err)
  }

  stdout, stderr := container.RunContainer(buff, "kardus:go")
  fmt.Println("STDOUT:" + stdout)
  fmt.Println("STDERR:" + stderr)
}
