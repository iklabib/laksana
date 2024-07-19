package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"codeberg.org/iklabib/laksana/util/fastrand"
)

func GetExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitError *exec.ExitError
	if ok := errors.As(err, &exitError); ok {
		return exitError.ExitCode()
	}

	return 0
}

func GetSignal(err error) (os.Signal, bool) {
	var exitError *exec.ExitError
	if ok := errors.As(err, &exitError); !ok {
		return nil, false
	}

	wt := exitError.Sys().(syscall.WaitStatus)
	if wt.Signaled() {
		return wt.Signal(), true
	}

	return nil, false
}

func signalAsMessage(signal os.Signal) string {
	switch signal {
	case syscall.SIGKILL:
		return "killed"
	case syscall.SIGXCPU:
		return "timeout"
	case syscall.SIGSEGV:
		return "out of memory"
	default:
		return "unknown reason"
	}
}

func ExitMessage(err error) string {
	signal, ok := GetSignal(err)
	if ok {
		return signalAsMessage(signal)
	}

	exitCode := GetExitCode(err)

	// FIXME: signal may not detected at times
	// subtracting by 128 and check
	switch exitCode - 128 {
	case 9:
		return signalAsMessage(syscall.SIGKILL)
	case 24:
		return signalAsMessage(syscall.SIGXCPU)
	case 11:
		return signalAsMessage(syscall.SIGSEGV)
	}

	// exit code 1 is likely because of failed test
	if exitCode == 0 || exitCode == 1 {
		return ""
	} else {
		return "internal error"
	}
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
