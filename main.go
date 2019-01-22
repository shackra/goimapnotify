package main

// Execute scripts on events using IDLE imap command (Go version)
// Copyright (C) 2017  Jorge Javier Araya Navarro

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
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	TLS        bool   `json:"tls,omitempty"`
	TLSOptions struct {
		RejectUnauthorized bool
	} `json:"tlsOption"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	PasswordCMD   string   `json:"passwordCmd,omitempty"`
	OnNewMail     string   `json:"onNewMail"`
	OnNewMailPost string   `json:"onNewMailPost,omitempty"`
	Boxes         []string `json:"boxes"`
}

// IDLEEvent models an IDLE event
type IDLEEvent struct {
	Mailbox   string
	EventType string
}

func main() {
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	fileconf := flag.String("conf", "path/to/imapnotify.conf", "Configuration file")
	list := flag.Bool("list", false, "List all mailboxes and exit")
	flag.Parse()
	raw, err := ioutil.ReadFile(*fileconf)
	if err != nil {
		log.Fatalf("[ERR] Can't read file: %s", err)
	}
	var conf NotifyConfig
	err = json.Unmarshal(raw, &conf)
	if err != nil {
		log.Fatalf("Can't parse the configuration: %s", err)
	}

	if *list {
		client := newClient(conf)
		defer client.Logout(30 * time.Second)
		_ = walkMailbox(client, "", 0)
	} else {
		events := make(chan IDLEEvent, 1)
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		guard := guardian{
			mx: &sync.Mutex{},
			wg: &sync.WaitGroup{},
		}

		NewWatchMailBox(conf, events, quit, &guard)

		// Send fake, first run, event
		events <- IDLEEvent{EventType: "EXISTS", Mailbox: "INBOX"}

		// Process incoming events from the mailboxes
		for rsp := range events {
			if rsp.EventType == "EXPUNGE" || rsp.EventType == "EXISTS" || rsp.EventType == "RECENT" {
				cmd := PrepareCommand(conf.OnNewMail, rsp)
				err := cmd.Run()
				if err != nil {
					log.Printf("[ERR] OnNewMail command failed: %s", err)
				} else {
					// execute the post command thing
					cmd := PrepareCommand(conf.OnNewMailPost, rsp)
					err := cmd.Run()
					if err != nil {
						log.Printf("[WARN] OnNewMailPost failed: %s", err)
					}
				}
			}
		}

		log.Println("[INF] Waiting for goroutines to finish...")
		guard.Wait()
		log.Println("Bye")
	}
}
