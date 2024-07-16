package main

// Execute scripts on events using IDLE imap command (Go version)
// Copyright (C) 2017-2024  Jorge Javier Araya Navarro

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
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	msg := donateMessage(8)
	fmt.Fprint(flag.CommandLine.Output(), "\n"+msg)
}

func main() {
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	fileconf := flag.String(
		"conf",
		filepath.Join(getDefaultConfigPath(), fmt.Sprintf("goimapnotify.%s", viper.SupportedExts[2])),
		"Configuration file",
	)
	list := flag.Bool("list", false, "List all mailboxes and exit")
	debug := flag.Bool("debug", false, "Output all network activity to the terminal")
	wait := flag.Int("wait", 1, "Period in seconds between IDLE event and execution of scripts")

	flag.Usage = usage

	flag.Parse()

	logrus.SetLevel(logrus.InfoLevel)
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infof("ℹ Running commit %s, tag %s, branch %s", commit, gittag, branch)

	viper.SetConfigFile(*fileconf)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Can't read file: '%s', error: %v", *fileconf, err)
	}

	topConfig, err := loadConfiguration(*fileconf)
	if err != nil {
		logrus.Fatalf("can't load the configuration: %v", err)
	}
	logrus.Debugf("configuration loaded successfuly: %s", *fileconf)

	idleChan := make(chan IDLEEvent)
	boxChan := make(chan BoxEvent, 1)
	quit := make(chan os.Signal, 1)
	doneChan := make(chan struct{})
	running := NewRunningBox(*debug, *wait)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	// indexes used because we need to change struct data
	for i := range topConfig.Configurations {
		topConfig.Configurations[i].Debug = *debug
		topConfig.Configurations[i] = retrieveCmd(topConfig.Configurations[i])
		if topConfig.Configurations[i].Alias == "" {
			if *debug {
				topConfig.Configurations[i].Alias = censorEmailAddress(
					topConfig.Configurations[i].Username,
				)
			} else {
				topConfig.Configurations[i].Alias = topConfig.Configurations[i].Username
			}
		}

		if *list {
			client, cErr := newClient(topConfig.Configurations[i])
			if cErr != nil {
				logrus.Fatalf("[%s] Something went wrong creating IMAP client: %s",
					topConfig.Configurations[i].Alias, cErr)
			}
			// nolint
			defer client.Logout()

			printDelimiter(client)
			_ = walkMailbox(client, "", 0)
		} else {
			// launch watchers for all mailboxes
			// listen in "boxes"

			for j := range topConfig.Configurations[i].Boxes {
				/*
				 * Copy default names if empty. Use SKIP to skip execution
				 * The check is happening in running.go:run
				 */
				topConfig.Configurations[i].Boxes[j] = setFromConfig(topConfig.Configurations[i], topConfig.Configurations[i].Boxes[j])
				client, iErr := newIMAPIDLEClient(topConfig.Configurations[i])
				if iErr != nil {
					logrus.Fatalf("[%s:%s] Something went wrong creating IDLE client: %s",
						topConfig.Configurations[i].Boxes[j].Alias, topConfig.Configurations[i].Boxes[j].Mailbox, iErr)
				}
				box := topConfig.Configurations[i].Boxes[j]
				key := box.Alias + box.Mailbox
				running.mutex[key] = new(sync.RWMutex)
				NewWatchBox(client, topConfig.Configurations[i], topConfig.Configurations[i].Boxes[j],
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
	printDonate(os.Stdout, 11)
	logrus.Info("Bye")
}
