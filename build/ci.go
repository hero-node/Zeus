package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GOBIN, _ = filepath.Abs(filepath.Join("build", "bin"))

func main() {
	doInstall()
}

func goPath() string {
	if os.Getenv("GOPATH") == "" {
		log.Fatal("No go env")
	}
	return os.Getenv("GOPATH")
}

var noGit bool

func gitcommit() string {
	cmd := exec.Command("git", []string{"log", "--format=\"%H\"", "-n", "1"}...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err == exec.ErrNotFound {
		if !noGit {
			noGit = true
			fmt.Println("NO git")
			return ""
		}
	} else if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSpace(stdout.String())
}

func doInstall() {
	//	gitCommit := gitcommit()
	//	gitCommitString := "main.gitCommit=" + gitCommit
	cmd := exec.Command("go", []string{"install", "--ldflags", "-v", "./gher.go"}...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOPATH="+goPath())
	cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
	cmd.Run()
}
