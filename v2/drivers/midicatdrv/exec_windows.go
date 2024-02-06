//go:build windows
// +build windows

package midicatdrv

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

/*
func execCommand(c string) *exec.Cmd {
	//return exec.Command("powershell.exe", "/Command",  `$Process = [Diagnostics.Process]::Start("` + c + `") ; echo $Process.Id `)
	//return exec.Command("powershell.exe", "/Command", `$Process = [Diagnostics.Process]::Start("fluidsynth.exe", "-i -q -n $_file") ; echo $Process.Id `)
	fmt.Println(c)
	return exec.Command("lib.exe", "/C", c)
}
*/

func midiCatOutCmd(index int) *exec.Cmd {
	cmd := exec.Command("midicat.exe", "out", fmt.Sprintf("--index=%v", index))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		// CREATE_NEW_CONSOLE
	}
	return cmd
}

func midiCatVersionCmd() *exec.Cmd {
	return exec.Command("midicat.exe", "version", "-s")
}

func midiCatInCmd(index int) *exec.Cmd {
	cmd := exec.Command("midicat.exe", "in", fmt.Sprintf("--index=%v", index))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	return cmd
}

func midiCatCmd(args string) *exec.Cmd {
	//return execCommand("midicat.exe " + args)
	//fmt.Println("midicat.exe " + args)
	a := strings.Split(args, " ")
	cmd := exec.Command("midicat.exe", a...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	return cmd
}
