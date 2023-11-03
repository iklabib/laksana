package main

// "instance" received a compiled program and run it.
// The program must be complete within 10 seconds, else be killed.
// "instance" supposed to be the entry point for container

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/iklabib/markisa/model"
	"gitlab.com/iklabib/markisa/util"
)

func main() {
	resp := model.RunResult{}
	buff, err := io.ReadAll(os.Stdin)
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
	}

	tempDir, err := os.MkdirTemp("", "box_")
	if err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
	}

	prog := filepath.Join(tempDir, "prog")
	if err := os.WriteFile(prog, buff, 0755); err != nil {
		resp.ExitCode = -1
		resp.Status = "INTERNAL_ERROR"
		jsonified, _ := json.Marshal(resp)
		fmt.Print(string(jsonified))
	}

	timeLimit := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.CommandContext(ctx, "./prog")
	cmd.Dir = tempDir

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	exitCode := 0
	status := "SUCCESS"
	if err := cmd.Run(); err != nil {
		switch ctx.Err() {
		case context.Canceled:
			status = "CANCELED"
		case context.DeadlineExceeded:
			status = "TIMEOUT"
		default:
			status = "ERROR"
		}
		exitCode = util.GetExitCode(&err)
		stderr.WriteString(err.Error())
	}

	resp.ExitCode = exitCode
	resp.Status = status
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()

	jsonified, _ := json.Marshal(resp)

	fmt.Print(string(jsonified))
}
