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
	wait := flag.Int("wait", 1, "Period in seconds between IDLE event and execution of scripts")

	flag.Parse()

	raw, err := ioutil.ReadFile(*fileconf)
	if err != nil {
		logrus.WithError(err).Fatalln("Can't read file")
	}

	idleChan := make(chan IDLEEvent)
	boxChan := make(chan BoxEvent, 1)
	quit := make(chan os.Signal, 1)
	doneChan := make(chan struct{})
	running := NewRunningBox(*debug, *wait)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	config, err := loadConfig(raw, *debug)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot read configuration")
	}

	// indexes used because we need to change struct data
	for i := range config {
		l := logrus.WithField("alias", config[i].Alias)
		if *list {
			client, cErr := newClient(config[i])
			if cErr != nil {
				l.WithError(cErr).Fatal("Something went wrong creating IMAP client")
			}
			// nolint
			defer client.Logout()

			_ = printDelimiter(client)
			_ = walkMailbox(client, "", 0)
		} else {
			// launch watchers for all mailboxes
			// listen in "boxes"

			for j := range config[i].Boxes {
				l = l.WithField("mailbox", config[i].Boxes[j].Mailbox)

				client, iErr := newIMAPIDLEClient(config[i])
				if iErr != nil {
					l.WithError(iErr).Fatal("Something went wrong creating IDLE client")
				}
				box := config[i].Boxes[j]
				key := box.Alias + box.Mailbox
				running.mutex[key] = new(sync.RWMutex)
				NewWatchBox(client, config[i], box, idleChan, boxChan, doneChan, wg)
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
	logrus.Info("Bye")
}
