package main

// Execute scripts on events using IDLE imap command (Go version)
// Copyright (C) 2017-2023  Jorge Javier Araya Navarro

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	commit string
	gittag string
	branch string
)

func getDefaultConfigPath() string {
	home := os.Getenv("XDG_CONFIG_HOME")
	if home == "" {
		return filepath.Join(os.Getenv("HOME"), ".config", "goimapnotify")
	}

	return filepath.Join(home, "goimapnotify")
}

func main() {
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	fileconf := flag.String("conf", filepath.Join(getDefaultConfigPath(), "goimapnotify.conf"), "Configuration file")
	list := flag.Bool("list", false, "List all mailboxes and exit")
	debug := flag.Bool("debug", false, "Output all network activity to the terminal (!! this may leak passwords !!)")
	wait := flag.Int("wait", 1, "Period in seconds between IDLE event and execution of scripts")

	flag.Parse()

	logrus.Infof("ℹ Running commit %s, tag %s, branch %s", commit, gittag, branch)

	raw, err := ioutil.ReadFile(*fileconf)
	if err != nil {
		logrus.Fatalf("Can't read file: %s", err)
	}
	var config []NotifyConfig
	err = json.Unmarshal(raw, &config)
	if err != nil {
		var configLegacy NotifyConfigLegacy
		err = json.Unmarshal(raw, &configLegacy)
		if err != nil {
			logrus.Fatalf("Can't parse the configuration: %s", err)
		} else {
			logrus.Warnf("Legacy configuration format detected")
			config = legacyConverter(configLegacy)
		}
	}

	idleChan := make(chan IDLEEvent)
	boxChan := make(chan BoxEvent, 1)
	quit := make(chan os.Signal, 1)
	doneChan := make(chan struct{})
	running := NewRunningBox(*debug, *wait)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	// indexes used because we need to change struct data
	for i := range config {
		config[i] = retrieveCmd(config[i])
		config[i].Debug = *debug
		if config[i].Alias == "" {
			config[i].Alias = config[i].Username
		}

		if *list {
			client, cErr := newClient(config[i])
			if cErr != nil {
				logrus.Fatalf("[%s] Something went wrong creating IMAP client: %s",
					config[i].Alias, cErr)
			}
			// nolint
			defer client.Logout()

			printDelimiter(client)
			_ = walkMailbox(client, "", 0)
		} else {
			// launch watchers for all mailboxes
			// listen in "boxes"

			for j := range config[i].Boxes {
				/*
				 * Copy default names if empty. Use SKIP to skip execution
				 * The check is happening in running.go:run
				 */
				config[i].Boxes[j] = setFromConfig(config[i], config[i].Boxes[j])
				client, iErr := newIMAPIDLEClient(config[i])
				if iErr != nil {
					logrus.Fatalf("[%s:%s] Something went wrong creating IDLE client: %s",
						config[i].Boxes[j].Alias, config[i].Boxes[j].Mailbox, iErr)
				}
				box := config[i].Boxes[j]
				key := box.Alias + box.Mailbox
				running.mutex[key] = new(sync.RWMutex)
				NewWatchBox(client, config[i], config[i].Boxes[j],
					idleChan, boxChan, doneChan, wg)
			}
		}
	}
	run := true
	if *list {
		run = false
	}
	for run {
		select {
		case boxEvent := <-boxChan:
			logrus.Infof("[%s:%s] Restarting watcher for mailbox",
				boxEvent.Mailbox.Alias, boxEvent.Mailbox.Mailbox)
			client, fErr := newIMAPIDLEClient(boxEvent.Conf)
			if fErr != nil {
				logrus.Fatalf("[%s:%s] Something went wrong creating IDLE client: %s",
					boxEvent.Mailbox.Alias, boxEvent.Mailbox.Mailbox, fErr)
			}
			NewWatchBox(client, boxEvent.Conf, boxEvent.Mailbox,
				idleChan, boxChan, doneChan, wg)
		case <-quit:
			// OS asked nicely to close, we ask our
			// goroutines to do the same
			close(doneChan)
			run = false
		case idleEvent := <-idleChan:
			wg.Add(1)
			go func() {
				defer wg.Done()
				running.schedule(idleEvent, doneChan)
			}()
		}
	}
	logrus.Info("Waiting other goroutines to stop...")
	wg.Wait()
	logrus.Info("Bye")
}
