package main

// "instance" received a compiled program and run it.
// The program must be complete within 10 seconds, else be killed.
// "instance" supposed to be the entry point for container

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// io.ReadAll(os.Stdin)
	buff, err := os.ReadFile("./hello-world.exe")
	if err != nil {
		panic(err)
	}

	tempDir, err := os.MkdirTemp("", "box_")
	if err != nil {
		panic(err)
	}
	prog := tempDir + "\\prog.exe"
	if err := os.WriteFile(prog, buff, 0755); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	var stdout strings.Builder
	var stderr strings.Builder

	status := "DONE"

	cmd := exec.CommandContext(ctx, ".\\prog.exe")
	cmd.Dir = tempDir

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}

	resp := Response{
		Status: status,
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	jsonified, _ := json.Marshal(resp)
	// compressed := util.Compress(jsonified)
	fmt.Println(string(jsonified))
}
