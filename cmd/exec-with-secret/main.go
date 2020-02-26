package main

import (
	"fmt"
	_ "github.com/nomeaning777/exec-with-secret/auto"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <command> <arguments>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	binary, err := exec.LookPath(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err := syscall.Exec(binary, os.Args[1:], os.Environ()); err != nil {
		panic(err)
	}
}
