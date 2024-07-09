package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"codeberg.org/iklabib/laksana/util/fastrand"
)

func GetExitCode(err *error) int {
	if (*err) == nil {
		return 0
	}

	if exitError, ok := (*err).(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	return 0
}

func CreateROFile(dest, data string) error {
	err := os.WriteFile(dest, []byte(data), 0o0444)
	if err != nil {
		return fmt.Errorf("failed to write to file")
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func RandomString() string {
	return fmt.Sprintf("%d", fastrand.Uint32())
}

func RandomBoxName() string {
	return fmt.Sprintf("box_%d", fastrand.Uint32())
}
