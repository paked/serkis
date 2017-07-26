package serkis

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
)

type Git struct {
	URL         string
	Username    string
	Password    string
	AuthorName  string
	AuthorEmail string
}

func (g *Git) PushNewChanges(root, fpath, message string) {
	err := g.Add(root, fpath)
	if err != nil {
		fmt.Println("FAILED. Could not add changes:", err)
		return
	}

	err = g.Commit(root, message)
	if err != nil {
		fmt.Println("FAILED. Could not commit changes:", err)
		return
	}

	err = g.Push(root, "origin", "master")
	if err != nil {
		fmt.Println("FAILED. Could not push changes:", err)
		return
	}
}

func (g *Git) PullRemoteChanges(root string) {
	err := g.Cmd(root, "git", "stash")
	if err != nil {
		fmt.Println("FAILED. Could not stash changes:", err)
		return
	}

	err = g.Cmd(root, "git", "pull", "--rebase")
	if err != nil {
		fmt.Println("FAILED. Could not rebase pull changes:", err)
		return
	}

	err = g.Cmd(root, "git", "push")
	if err != nil {
		fmt.Println("FAILED. Could not push changes:", err)
		return
	}

	err = g.Cmd(root, "git", "stash", "apply", "--index")
	if err != nil {
		fmt.Println("WARNING. Could not apply stashed changes:", err)
		return
	}
}

func (g *Git) Clone(root string) error {
	gURL, err := url.Parse(g.URL)
	if err != nil {
		return err
	}

	gURL.User = url.UserPassword(g.Username, g.Password)

	return g.Cmd(".", "git", "clone", gURL.String(), root)
}

func (g *Git) Add(root, path string) error {
	return g.Cmd(root, "git", "add", path)
}

func (g *Git) Commit(root, message string) error {
	return g.Cmd(root, "git", "commit", "--message", message)
}

func (g *Git) Push(root, remote, branch string) error {
	return g.Cmd(root, "git", "push", remote, branch)
}

func (g *Git) Pull() {
	// git pull
}

func (g *Git) Config(root, field, value string) error {
	return g.Cmd(root, "git", "config", field, value)
}

func (g *Git) Author() string {
	return fmt.Sprintf("%s <%s>", g.AuthorName, g.AuthorEmail)
}

func (g *Git) Cmd(root, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
