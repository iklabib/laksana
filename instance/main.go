package main

// "instance" received a compiled program and run it.
// The program must be complete within 10 seconds, else be killed.
// "instance" supposed to be the entry point for container

import (
	"context"
	"encoding/json"
	"markisa/model"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {
	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	tempDir, err := os.MkdirTemp("", "box_")
	if err != nil {
		panic(err)
	}

	prog := path.Join(tempDir, "prog")
	if err := os.WriteFile(prog, buff, 0755); err != nil {
		panic(err)
	}
  
  timeLimit := time.Second*10
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.CommandContext(ctx, "./prog")
	cmd.Dir = tempDir

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

  status := "SUCCESS"
  if err := cmd.Run(); err != nil {
    status = "ERROR"
    stderr.WriteString(err.Error())
  } else if err := ctx.Err(); err != nil {
    switch err {
      case context.Canceled:
          status = "CANCELED"
      case context.DeadlineExceeded:
          status = "TIMEOUT"
      default:
          status = "ERROR"
    }
    stderr.WriteString(err.Error())
  } 


  resp := model.Response{
		Status: status,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	jsonified, _ := json.Marshal(resp)

  // send result to stdout and catch it outside of container
  fmt.Print(string(jsonified))
}
