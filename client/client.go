package client

import (
	"time"

	"github.com/ao-data/albiondata-client/log"
)

var version string

// Client struct base
type Client struct {
}

// NewClient return a new Client instance
func NewClient(_version string) *Client {
	version = _version
	return &Client{}
}

// Run starts client settings and run
func (client *Client) Run() error {
	log.Infof("Starting Albion Data Client, version: %s", version)
	log.Info("This is a third-party application and is in no way affiliated with Sandbox Interactive or Albion Online.")
	log.Info("Additional parameters can listed by calling this file with the -h parameter.")

	ConfigGlobal.setupDebugEvents()
	ConfigGlobal.setupDebugOperations()

	createDispatcher()

	if ConfigGlobal.Offline {
		processOffline(ConfigGlobal.OfflinePath)

		// TODO: get rid of this, some kind of stupid delay locally when hitting /pow on local dev server, fix it
		n := 1
		for n < 5000 {
			time.Sleep(1 * time.Millisecond)
			n++
		}

	} else {
		apw := newAlbionProcessWatcher()
		return apw.run()
	}
	return nil
}
