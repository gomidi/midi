//go:build !windows
// +build !windows

package midicatdrv

import (
	"fmt"
	"os/exec"
	"syscall"
)

func _execCommand(c string) *exec.Cmd {
	return exec.Command("sh", "-c", "exec "+c)
}

func midiCatOutCmd(index int) *exec.Cmd {
	cmd := exec.Command("midicat", "out", fmt.Sprintf("--index=%v", index))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	return cmd
}

func midiCatVersionCmd() *exec.Cmd {
	return exec.Command("midicat", "version", "-s")
}

func midiCatInCmd(index int) *exec.Cmd {
	cmd := exec.Command("midicat", "in", fmt.Sprintf("--index=%v", index))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	return cmd
}

func midiCatCmd(args string) *exec.Cmd {
	cmd := _execCommand("midicat " + args)
	// important! prevents that signals such as interrupt send to the main program gets passed
	// to midicat (which would not allow us to shutdown properly, e.g. stop hanging notes)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}
	return cmd
}
