package zig

import (
	"markisa/container"
	"markisa/util"
	"os"
	"os/exec"
	"path"
)

func Run(archive string) (string, string, error){
  tempDir, _ := util.CreateTempDir()
  bin, err := Build(archive, tempDir)
  if (err != nil) {
    return "", "", err
  }

  stdout, stderr, err := container.RunContainer(bin, "localhost/markisa:zig")
  if err != nil {
    return "", "", err
  }
  return stdout, stderr, err
}

func Build(archive string, dir string) ([]byte, error) {
  srcPath := path.Join(dir, "prog.zig")
  src, _ := os.Create(srcPath)
  defer src.Close()
  if _, err := src.WriteString(archive); err != nil {
    panic(err)
  }

  // TODO: incremental build
  cmd := exec.Command("zig", "build-exe", "prog.zig")
  cmd.Dir = dir
  if err := cmd.Run(); err != nil {
    return nil, err
  }

  // TODO: return compile error stage
  prog, _ := os.ReadFile(path.Join(dir, "prog"))
  os.RemoveAll(dir)

  return prog, nil
}

func initProject(dir string) error {
  cmd := exec.Command("zig", "init-exe")
  cmd.Dir = dir
  err := cmd.Run()
  return err
}
