package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GOBIN, _ = filepath.Abs(filepath.Join("build", "bin"))

func main() {
	switch os.Args[1] {
	case "install":
		doInstall()
	case "xgo":
		doXgo(os.Args[2:])
	}
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

	cmd := exec.Command("go", []string{"install", "-v", "./cmd/gher/gher.go"}...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	cmd.Env = append(cmd.Env, "GOPATH="+goPath())
	cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)

	cmd.Run()
}

func getXgo() {
	cmd := exec.Command("go", []string{"get", "github.com/karalabe/xgo"}...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	cmd.Env = append(cmd.Env, "GOPATH="+goPath())
	cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)

	cmd.Run()
}

func doXgo(cmdlinne []string) {
	flag.CommandLine.Parse(cmdlinne)
	getXgo()

	args := append([]string{}, flag.Args()...)

	path := "./cmd/gher"
	args = append(args, []string{"--dest", GOBIN, path}...)
	cmd := exec.Command(filepath.Join(GOBIN, "xgo"), args...)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GOPATH=") || strings.HasPrefix(e, "GOBIN=") {
			continue
		}
		cmd.Env = append(cmd.Env, e)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, "GOPATH="+goPath())
	cmd.Env = append(cmd.Env, "GOBIN="+GOBIN)
	fmt.Println(cmd.Args)
	cmd.Run()
}
