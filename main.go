package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/teejays/clog"
)

const (
	FILE_NAME string = `changeme.txt`

	DEFAULT_NUM_ACTIVITIES_MIN int = 0
	DEFAULT_NUM_ACTIVITIES_MAX int = 6

	DEFAULT_WAIT_DURATION_SECS_MIN int = 0
	DEFAULT_WAIT_DURATION_SECS_MAX int = 3600

	MAX_COMMIT_MESSAGE_LEN int = 200
)

var (
	ErrCommitMessageTooLong error = fmt.Errorf("Commit message is too long, which could mean an error.")
)

var (
	gitLock  sync.Mutex
	fileLock sync.Mutex
)

var (
	dev             bool
	dirPath         string
	numActivityMin  int
	numActivityMax  int
	waitDurationMin int
	waitDurationMax int
)

func main() {
	clog.Info("Initializing the githubber...")

	clog.LogToSyslog = true

	flag.BoolVar(&dev, "dev", false, "Pass this flag when running in development mode to skip commits and pushes.")
	flag.StringVar(&dirPath, "dir", "", "The path for the git repo, where there is a file named changeme.txt.")

	flag.IntVar(&numActivityMin, "min", DEFAULT_NUM_ACTIVITIES_MIN, "Pass this to set the minimum number of activities allowed.")
	flag.IntVar(&numActivityMax, "max", DEFAULT_NUM_ACTIVITIES_MAX, "Pass this to set the maximum number of activities allowed.")

	flag.IntVar(&waitDurationMin, "wait-min", DEFAULT_WAIT_DURATION_SECS_MIN, "Pass this to override the minimum number of seconds to wait before an activity.")
	flag.IntVar(&waitDurationMax, "wait-max", DEFAULT_WAIT_DURATION_SECS_MAX, "Pass this to override the maximum number of seconds to wait before an activity.")

	flag.Parse()

	if strings.TrimSpace(dirPath) == "" {
		clog.Fatal("Invalid dirPath: cannot be empty")
	}
	if numActivityMin < 0 {
		clog.Fatal("Invalid numActivityMin: cannot be less than 0")
	}
	if numActivityMax < 0 {
		clog.Fatal("Invalid numActivityMax: cannot be less than 0")
	}
	if waitDurationMin < 0 {
		clog.Fatal("Invalid waitDurationMin: cannot be less than 0")
	}
	if waitDurationMax < 0 {
		clog.Fatal("Invalid waitDurationMax: cannot be less than 0")
	}

	if dev {
		clog.Notice("Running in DEVELOPMENT mode")
	}

	// How many commits should I make?
	var numActivities int = getRandomInt(numActivityMin, numActivityMax)
	clog.Infof("Number of activities to be made right now: %d", numActivities)

	// WaitGroup for concurrency
	var wg sync.WaitGroup

	for i := 0; i < numActivities; i++ {

		// For each activity, we should probably add some randomized wait times (between 1 and 3600 secs)
		// so the commit history looks natural,
		waitDuration := time.Second * time.Duration(getRandomInt(waitDurationMin, waitDurationMax))

		// Add a counter to wait group so we can keep track of how many concurrent goroutines are running
		wg.Add(1)

		// Each activity is going to be it's own goroutine.
		go func(i int, wait time.Duration) {

			// Catch any panics that happen in this goroutine
			defer func() {
				if r := recover(); r != nil {
					clog.Errorf("Panic in goroutine (recovered): %s", r)
				}
				wg.Done()
			}()

			clog.Infof("[Activity %d / %d] Going to wait for %s before doing the activity", i+1, numActivities, wait)

			// Wait for sometime before doing the activity
			time.Sleep(wait)

			clog.Infof("[Activity %d / %d] Processing activity...", i+1, numActivities)

			// Do the activity
			err := doActivity()
			if err != nil {
				clog.FatalErr(err)
			}

			clog.Infof("[Activity %d / %d] Finished activity...", i+1, numActivities)

		}(i, waitDuration)
	}

	// Main goroutine needs to wait for all the sub-goroutines to finish
	wg.Wait()

	clog.Infof("Finished %d activities.", numActivities)

	return
}

func doActivity() (err error) {

	clog.Debug(" - Starting to code...")
	err = doCoding()
	if err != nil {
		return err
	}
	clog.Debug(" - Finished coding.")

	clog.Debug(" - Starting to add...")
	err = doGitAdd()
	if err != nil {
		return err
	}
	clog.Debug(" - Finished adding.")

	clog.Debug(" - Starting to commit...")
	err = doGitCommit()
	if err != nil {
		return err
	}
	clog.Debug(" - Finished commiting.")

	clog.Debug(" - Starting to push...")
	err = doGitPush()
	if err != nil {
		return err
	}
	clog.Debug(" - Finished pushing.")

	return

}

func doCoding() error {
	fileLock.Lock()
	defer fileLock.Unlock()

	filePath := filepath.Join(dirPath, FILE_NAME)

	// read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// change/append something
	newText := fmt.Sprintf("This is a change made on %s.\n", time.Now().Format(time.RFC1123Z))
	newContent := append(content, []byte(newText)...)
	// Write the file
	err = ioutil.WriteFile(filePath, newContent, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

func doGitAdd() error {
	cmd := exec.Command("git", "-C", dirPath, "add", FILE_NAME)
	var out bytes.Buffer
	cmd.Stdout = &out

	gitLock.Lock()
	defer gitLock.Unlock()

	err := cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}
	clog.Debug(out.String())
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

	cmd := exec.Command("git", "-C", dirPath, "commit", "-m", commitMessage)
	var out bytes.Buffer
	cmd.Stdout = &out

	if dev {
		clog.Warnf("Commits disabled during development mode. Not commiting any changes.")
		return nil
	}

	gitLock.Lock()
	defer gitLock.Unlock()

	err = cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}

	clog.Debug(out.String())

	return nil
}

func doGitPush() error {
	cmd := exec.Command("git", "-C", dirPath, "push")
	var out bytes.Buffer
	cmd.Stdout = &out

	if dev {
		clog.Warnf("Pushed disabled during development mode. Not pushing any changes.")
		return nil
	}

	gitLock.Lock()
	defer gitLock.Unlock()

	err := cmd.Run()
	if err != nil {
		clog.Error(out.String())
		return err
	}

	clog.Debug(out.String())

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

	if len(body) > MAX_COMMIT_MESSAGE_LEN {
		clog.Warnf("Commit message is too long:\n%s", string(body))
		return "", ErrCommitMessageTooLong
	}
	return string(body), nil
}

func getDefaultCommitMessage() string {
	return fmt.Sprintf("This is a test commit on %s", time.Now().Format(`Monday, Jan 2 2006 3:04PM`))
}
