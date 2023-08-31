package c

import (
	"markisa/container"
	"markisa/model"
	"markisa/util"
	"os"
	"os/exec"
	"path"
	"strings"
)

func Run(archive string) model.RunResponse {
  tempDir, _ := util.CreateTempDir()
  bin, buildResult := Build(archive, tempDir)

  runResponse := model.RunResponse {
    Build: buildResult,
  }

  if (buildResult.ExitCode != 0) {
    return runResponse
  }

  runResponse.Run = container.RunContainer(bin, "localhost/markisa:common")
  return runResponse
}

func Build(archive string, dir string) ([]byte, model.BuildResult) {
  srcPath := path.Join(dir, "prog.c")
  src, _ := os.Create(srcPath)
  defer src.Close()
  src.WriteString(archive)

  var stdout strings.Builder
  var stderr strings.Builder

  // I just don't feel like to install clang or gcc
  cmd := exec.Command("zig", "cc", "prog.c", "-o", "prog")
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  cmd.Dir = dir

  err := cmd.Run(); 

  buildResult := model.BuildResult {
    ExitCode: util.GetExitCode(&err),
    Stdout: stdout.String(),
    Stderr: stderr.String(),
  }

  prog, _ := os.ReadFile(path.Join(dir, "./prog"))
  return prog, buildResult
}