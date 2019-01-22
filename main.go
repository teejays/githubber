package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/teejays/clog"
)

const (
	FILE_NAME string = `dummy.txt`
)

func main() {
	clog.Info("Initializing the githubber...")

	// How many commits should I make?
	numCommits := getRandomInt(1, 6)
	if numCommits > 4 { // reduce the likelihood of 6 commits
		numCommits = getRandomInt(3, 6)
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
	content, err := ioutil.ReadFile(FILE_NAME)
	if err != nil {
		return err
	}

	// change/append something
	newText := fmt.Sprintf("This is a test commit on %s.\n", time.Now().Format(time.RFC1123Z))
	newContent := append(content, []byte(newText)...)
	// Write the file
	err = ioutil.WriteFile(FILE_NAME, newContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

func doGitAdd() error {
	cmd := exec.Command("git", "add", FILE_NAME)
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
	commitMessage, err := getRandomCommitMessage()
	if err != nil {
		commitMessage = getDefaultCommitMessage()
		clog.Warnf("There was an error generating a commit message: %s\nUsing standard test commit message: %s", err, commitMessage)
		return nil
	}

	clog.Debugf("Using the commit message: %s", commitMessage)

	cmd := exec.Command("git", "commit", "-m", commitMessage)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
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

func getRandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func getRandomCommitMessage() (string, error) {
	resp, err := http.Get("http://whatthecommit.com/index.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getDefaultCommitMessage() string {
	return fmt.Sprintf("This is a test commit on %s", time.Now().Format(`Monday, Jan 2 2006 3:04PM`))
}
