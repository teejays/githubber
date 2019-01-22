package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/teejays/clog"
)

const fileName string = `dummy.txt`

func main() {
	var err error

	clog.Info("Starting to code...")
	err = doCoding()
	if err != nil {
		clog.FatalErr(err)
	}
	clog.Info("Finished coding.")

	clog.Info("Starting to add...")
	err = doGitAdd()
	if err != nil {
		clog.FatalErr(err)
	}
	clog.Info("Finished adding.")

	clog.Info("Starting to commit...")
	err = doGitCommit()
	if err != nil {
		clog.FatalErr(err)
	}
	clog.Info("Finished commiting.")

	clog.Info("Starting to push...")
	err = doGitPush()
	if err != nil {
		clog.FatalErr(err)
	}
	clog.Info("Finished pushing.")

	clog.Info("Finished!")
}

func doCoding() error {
	// read the file
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	// change/append something
	newText := "This is today's commit.\n"
	newContent := append(content, []byte(newText)...)

	// Write the file
	err = ioutil.WriteFile(fileName, newContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

func doGitAdd() error {
	cmd := exec.Command("git", "add", fileName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}
	clog.Info(out.String())
	return nil
}
func doGitCommit() error {
	cmd := exec.Command("git", "commit", "-m", "'Made editions to the the file'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}
	clog.Info(out.String())

	return nil
}

func doGitPush() error {
	cmd := exec.Command("git", "push")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}
	clog.Info(out.String())
	return nil
}
