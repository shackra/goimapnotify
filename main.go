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
	"slices"
	"sync"
	"syscall"

	"github.com/emersion/go-imap"

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
				logrus.Fatalf("[%s] Something went wrong creating IMAP client: %v",
					topConfig.Configurations[i].Alias, cErr)
			}
			// nolint
			defer client.Logout()

			max, err := printDelimiter(client)
			if err != nil {
				logrus.WithError(err).Warning("listing mailboxes finished with error")
			}
			_ = walkMailbox(client, "", 0, max)
		} else {
			if len(topConfig.Configurations[i].Boxes) == 0 {
				client, iErr := newIMAPIDLEClient(topConfig.Configurations[i])
				if iErr != nil {
					logrus.Fatalf(iErr.Error())
				}
				mailboxes := make(chan *imap.MailboxInfo, 10)
				done := make(chan error, 1)
				go func() {
					done <- client.List("", "*", mailboxes)
				}()
				var mboxes []string

				for m := range mailboxes {
					if slices.Contains([]string{"[Gmail]", "[Gmail]/All Mail"}, m.Name) {
						continue
					}
					mboxes = append(mboxes, m.Name)
				}
				for _, m := range mboxes {
					client, iErr := newIMAPIDLEClient(topConfig.Configurations[i])
					if iErr != nil {
						logrus.Fatalf(iErr.Error())
					}
					nConf := topConfig.Configurations[i]
					box := Box{
						Alias:             nConf.Alias,
						Mailbox:           m,
						OnNewMail:         nConf.OnNewMail,
						OnNewMailPost:     nConf.OnNewMailPost,
						OnDeletedMail:     nConf.OnDeletedMail,
						OnDeletedMailPost: nConf.OnDeletedMailPost,
					}
					key := box.Alias + box.Mailbox
					running.mutex[key] = new(sync.RWMutex)
					NewWatchBox(client, NotifyConfig{}, box, idleChan, boxChan, doneChan, wg)
				}
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
	}

	run := !*list
	for run {
		select {
		case boxEvent := <-boxChan:
			l := logrus.WithField("alias", boxEvent.Mailbox.Alias).WithField("mailbox", boxEvent.Mailbox.Mailbox)
			l.Info("Restarting watcher for mailbox")
			client, fErr := newIMAPIDLEClient(boxEvent.Conf)
			if fErr != nil {
				l.WithError(fErr).Fatalf("Something went wrong creating IDLE client")
			}
			NewWatchBox(client, boxEvent.Conf, boxEvent.Mailbox, idleChan, boxChan, doneChan, wg)
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
	printDonate(os.Stderr, 11)
	logrus.Info("Bye")
}
