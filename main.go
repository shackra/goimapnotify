package main

// Execute scripts on events using IDLE imap command (Go version)
// Copyright (C) 2017-2021  Jorge Javier Araya Navarro

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
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	fileconf := flag.String("conf", "path/to/imapnotify.conf", "Configuration file")
	list := flag.Bool("list", false, "List all mailboxes and exit")
	debug := flag.Bool("debug", false, "Output all network activity to the terminal (!! this may leak passwords !!)")

	flag.Parse()

	raw, err := ioutil.ReadFile(*fileconf)
	if err != nil {
		logrus.Fatalf("[ERR] Can't read file: %s", err)
	}
	var conf NotifyConfig
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		logrus.Fatalf("Can't parse the configuration: %s", err)
	}
	conf.Debug = *debug

	if *list {
		client, cErr := newClient(conf)
		if cErr != nil {
			logrus.Fatalf("something went wrong creating IMAP client: %s", cErr)
		}
		// nolint
		defer client.Logout()

		printDelimiter(client)
		_ = walkMailbox(client, "", 0)
	} else {
		events := make(chan IDLEEvent)
		boxEvents := make(chan BoxEvent, 1)
		quit := make(chan os.Signal, 1)
		done := make(chan struct{})
		runningBoxes := NewRunningBox(conf.Wait, conf.Debug)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		wg := &sync.WaitGroup{}

		// Send "Good Morning!" event
		boxEvents <- BoxEvent{Type: FirstStart}
		run := true
		for run {
			select {
			case box := <-boxEvents:
				if box.Type == FirstStart {
					// launch watchers for all mailboxes
					// listen in "boxes"
					for _, mailbox := range conf.Boxes {
						client, iErr := newIMAPIDLEClient(conf)
						if iErr != nil {
							logrus.Fatalf("something went wrong creating IDLE client: %s", iErr)
						}
						NewWatchBox(client, mailbox, events, boxEvents, done, wg)
					}
				} else {
					logrus.Infof("restarting watcher for mailbox %s", box.Mailbox)
					client, fErr := newIMAPIDLEClient(conf)
					if fErr != nil {
						logrus.Fatalf("something went wrong creating IDLE client: %s", fErr)
					}
					NewWatchBox(client, box.Mailbox, events, boxEvents, done, wg)
				}
			case <-quit:
				// OS asked nicely to close, we ask our
				// goroutines to do the same
				close(done)
				run = false
			case idle := <-events:
				runningBoxes.RunOrIgnore(conf.OnNewMail, conf.OnNewMailPost, idle)
			}
		}

		logrus.Info("waiting other goroutines to stop...")
		wg.Wait()
		logrus.Info("Bye")
	}
}
