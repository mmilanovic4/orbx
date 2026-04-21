package sysutil

import (
	"os"
	"path"
)

func GetShell() string {
	shell := os.Getenv("SHELL")
	return path.Base(shell)
}
