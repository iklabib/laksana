package c

import (
	"gitlab.com/iklabib/markisa/container"
	"gitlab.com/iklabib/markisa/model"
	"gitlab.com/iklabib/markisa/util"
	"os"
	"os/exec"
	"path/filepath"
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
  srcPath := filepath.Join(dir, "prog.c")
  src, _ := os.Create(srcPath)
  defer src.Close()
  src.WriteString(archive)

  var stdout strings.Builder
  var stderr strings.Builder

  cmd := exec.Command("gcc", "prog.c", "-o", "prog")
  cmd.Stdout = &stdout
  cmd.Stderr = &stderr
  cmd.Dir = dir

  err := cmd.Run(); 

  buildResult := model.BuildResult {
    ExitCode: util.GetExitCode(&err),
    Stdout: stdout.String(),
    Stderr: stderr.String(),
  }

  prog, _ := os.ReadFile(filepath.Join(dir, "./prog"))
  return prog, buildResult
}
