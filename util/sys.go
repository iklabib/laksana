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

// SIGIO, SIGIOT, SIGCLD, SIGUNUSED are archaic
var knownSignals map[syscall.Signal]bool = map[syscall.Signal]bool{
	syscall.SIGABRT:   true,
	syscall.SIGALRM:   true,
	syscall.SIGBUS:    true,
	syscall.SIGCHLD:   true,
	syscall.SIGCONT:   true,
	syscall.SIGFPE:    true,
	syscall.SIGHUP:    true,
	syscall.SIGILL:    true,
	syscall.SIGINT:    true,
	syscall.SIGIO:     true,
	syscall.SIGKILL:   true,
	syscall.SIGPIPE:   true,
	syscall.SIGPROF:   true,
	syscall.SIGPWR:    true,
	syscall.SIGQUIT:   true,
	syscall.SIGSEGV:   true,
	syscall.SIGSTKFLT: true,
	syscall.SIGSTOP:   true,
	syscall.SIGSYS:    true,
	syscall.SIGTERM:   true,
	syscall.SIGTRAP:   true,
	syscall.SIGTSTP:   true,
	syscall.SIGTTIN:   true,
	syscall.SIGTTOU:   true,
	syscall.SIGURG:    true,
	syscall.SIGUSR1:   true,
	syscall.SIGUSR2:   true,
	syscall.SIGVTALRM: true,
	syscall.SIGWINCH:  true,
	syscall.SIGXCPU:   true,
	syscall.SIGXFSZ:   true,
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

	// FIXME: signal may not detected
	// instead minijail return 128 + signal as exit code
	sig := syscall.Signal(GetExitCode(err) - 128)
	if knownSignals[sig] {
		return sig, true
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
