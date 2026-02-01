package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ao-data/albiondata-client/client"
	"github.com/ao-data/albiondata-client/log"
	"github.com/ao-data/albiondata-client/systray"

	"github.com/ao-data/go-githubupdate/updater"
)

var version string

func init() {
	client.ConfigGlobal.SetupFlags()
}

func main() {
	if client.ConfigGlobal.PrintVersion {
		log.Infof("Albion Data Client, version: %s", version)
		return
	}

	startUpdater()

	go systray.Run()

	c := client.NewClient(version)
	err := c.Run()
	if err != nil {
		log.Error(err)
		log.Error("The program encountered an error. Press any key to close this window.")
		var b = make([]byte, 1)
		_, _ = os.Stdin.Read(b)
	}

}

func startUpdater() {
	if version != "" && !strings.Contains(version, "dev") {
		u := updater.NewUpdater(
			version,
			client.ConfigGlobal.UpdateGithubOwner,
			client.ConfigGlobal.UpdateGithubRepo,
			"update-",
		)

		go func() {
			for {
				if tryUpdate(u) {
					restartProcess()
					return // This line won't be reached if restart succeeds, but included for clarity
				}
				// Wait 1 hour before checking again
				time.Sleep(time.Hour)
			}
		}()
	}
}

// tryUpdate attempts to check and apply an update with retry logic.
// Returns true if an update was successfully applied.
func tryUpdate(u *updater.Updater) bool {
	maxTries := 2
	for i := 0; i < maxTries; i++ {
		updated, err := u.BackgroundUpdater()
		if err != nil {
			log.Error(err.Error())
			if i < maxTries-1 {
				log.Info("Will try again in 60 seconds. You may need to run the client as Administrator.")
				time.Sleep(time.Second * 60)
			}
			continue
		}
		if updated {
			return true
		}
		// No update available, no need to retry
		return false
	}
	return false
}

// restartProcess starts a new instance of the application and exits the current process.
func restartProcess() {
	execPath, err := os.Executable()
	if err != nil {
		log.Errorf("Failed to get executable path for restart: %v", err)
		return
	}

	log.Info("Restarting with updated version...")

	// Start the new process with the same arguments
	cmd := exec.Command(execPath, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		log.Errorf("Failed to start new process: %v", err)
		return
	}

	log.Info("New process started, exiting current process.")
	os.Exit(0)
}
