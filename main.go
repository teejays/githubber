package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/teejays/clog"
)

const fileName string = `dummy.txt`

func main() {
	clog.Info("Initializing the githubber...")

	// How many commits should I make?
	numCommits := rand.Intn(6)
	if numCommits > 4 { // reduce the likelihood of 6 commits
		numCommits = rand.Intn(6)
	}

	clog.Infof("Number of commits to be made right now: %d", numCommits)

	for i := 0; i < numCommits; i++ {
		clog.Infof("Processing Commit %d of %d...", i+1, numCommits)
		err := doActivity()
		if err != nil {
			clog.FatalErr(err)
		}
	}

	clog.Infof("Finished %d commits.", numCommits)

	return
}

func doActivity() (err error) {

	clog.Info("Starting activity...")

	clog.Info(" - Starting to code...")
	err = doCoding()
	if err != nil {
		return err
	}
	clog.Info(" - Finished coding.")

	clog.Info(" - Starting to add...")
	err = doGitAdd()
	if err != nil {
		return err
	}
	clog.Info(" - Finished adding.")

	clog.Info(" - Starting to commit...")
	err = doGitCommit()
	if err != nil {
		return err
	}
	clog.Info(" - Finished commiting.")

	clog.Info(" - Starting to push...")
	err = doGitPush()
	if err != nil {
		return err
	}
	clog.Info(" - Finished pushing.")

	clog.Info("Finished activity.")

	return

}

func doCoding() error {
	// read the file
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	// change/append something
	newText := fmt.Sprintf("This is a test commit on %s.\n", time.Now().Format(time.RFC1123Z))
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
	cmd := exec.Command("git", "commit", "-m", fmt.Sprintf("'Made editions to the the file on %s'", time.Now().Format(time.RFC1123Z)))
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
