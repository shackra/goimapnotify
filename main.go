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
	"strings"
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
	loglevel := flag.String("log-level", "info", "change the logging level, possible values: error, warning/warn, info/information, debug")
	wait := flag.Int("wait", 1, "Period in seconds between IDLE event and execution of scripts")

	flag.Usage = usage

	flag.Parse()

	debug := false

	switch strings.ToLower(*loglevel) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		debug = true
	case "info", "information":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.Fatalf("unknown logging level %q", *loglevel)
	}

	logrus.Infof("ℹ Running commit %s, tag %s, branch %s", commit, gittag, branch)

	viper.SetConfigFile(*fileconf)
	if err := viper.ReadInConfig(); err != nil {
		logrus.WithError(err).Fatalf("can't read file: %q", *fileconf)
	}

	idleChan := make(chan IDLEEvent)
	boxChan := make(chan BoxEvent, 1)
	quit := make(chan os.Signal, 1)
	quitChan := make(chan struct{})

	topConfig, err := loadConfiguration(*fileconf)
	if err != nil {
		logrus.WithError(err).Fatalf("can't load the configuration %q", *fileconf)
	}
	logrus.Debugf("configuration loaded successfully: %q", *fileconf)

	running := NewRunningBox(debug, *wait)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	if *list {
		for _, account := range topConfig.Configurations {
			client, err := newClient(account)
			if err != nil {
				logrus.WithError(err).WithField("account", account.Alias).Fatal("something went wrong creating IMAP client")
			}
			// nolint
			defer client.Logout()

			max, err := printDelimiter(client)
			if err != nil {
				logrus.WithField("alias", account.Alias).WithError(err).Warning("listing mailboxes finished with error")
			}
			logrus.WithField("account", account.Alias).Info("walking through the account mailboxes")
			err = walkMailbox(client, "", 0, max)
			if err != nil {
				logrus.WithField("account", account.Alias).WithError(err).Fatal("something went wrong while walking on the account listing all mailboxes")
			}
		}
	}

	// Watch mailboxes events
	// This kick-starts the watching
	idleForever := !*list
	if idleForever {
		/* I really doubt it that creating a new client for
		   each mailbox that we want to listen for events is
		   healthy, or elegant... but, if the connection
		   fails, what the program does right now is exactly
		   that: it creates a new client for that failing
		   mailbox only, lol!
		*/
		for _, account := range topConfig.Configurations {
			running.mutex[account.Alias] = new(sync.RWMutex)
			for _, mailbox := range account.Boxes {
				client, err := newIMAPIDLEClient(account)
				if err != nil {
					logrus.WithError(err).WithField("account", account.Alias).Fatal("cannot make IMAP client")
				}
				key := account.Alias + mailbox.Mailbox
				running.mutex[key] = new(sync.RWMutex)
				running.config[key] = account
				NewWatchBox(client, account, mailbox, idleChan, boxChan, quitChan, wg)
			}
		}
	}

	for idleForever {
		select {
		case boxEvent := <-boxChan:
			key := boxEvent.Mailbox.Alias + boxEvent.Mailbox.Mailbox
			l := logrus.WithField("alias", boxEvent.Mailbox.Alias).WithField("mailbox", boxEvent.Mailbox.Mailbox)
			l.Info("Restarting watcher for mailbox")
			client, fErr := newIMAPIDLEClient(running.config[key])
			if fErr != nil {
				l.WithError(fErr).Fatal("Something went wrong creating IDLE client")
			}
			NewWatchBox(client, running.config[key], boxEvent.Mailbox, idleChan, boxChan, quitChan, wg)
		case <-quit:
			// OS asked nicely to close, we ask our
			// goroutines to do the same
			close(quitChan)
			idleForever = false
		case idleEvent := <-idleChan:
			wg.Add(1)
			go func() {
				defer wg.Done()
				running.schedule(idleEvent, quitChan)
			}()
		}
	}
	logrus.Info("waiting other goroutines to stop...")
	wg.Wait()
	printDonate(os.Stderr, 11)
	logrus.Info("bye")
}
