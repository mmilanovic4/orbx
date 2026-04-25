package sysutil

import (
	"os"
	"os/exec"
	"path"
)

func GetShell() string {
	shell := os.Getenv("SHELL")
	return path.Base(shell)
}

func ClearScreen() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
