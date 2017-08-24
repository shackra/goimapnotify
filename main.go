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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mxk/go-imap/imap"
)

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host       string
	Port       int  `json:",omitempty"`
	TLS        bool `json:",omitempty"`
	TLSOptions struct {
		RejectUnauthorized bool
	} `json:",omitempty"`
	Username      string
	Password      string
	OnNewMail     string
	OnNewMailPost string
	Boxes         []string
}

// IDLEEvent models an IDLE event
type IDLEEvent struct {
	Mailbox   string
	EventType string
}

// PrepareCommand parse a string and return a command executable by Go
func PrepareCommand(command string, rsp IDLEEvent) *exec.Cmd {
	var commandstr string
	if strings.Contains("%s", command) {
		commandstr = fmt.Sprintf(command, rsp.Mailbox)
	} else {
		commandstr = command
	}
	commandsplt := strings.Split(commandstr, " ")
	commandhead := commandsplt[0]
	args := commandsplt[:1]
	cmd := exec.Command(commandhead, args...)
	return cmd
}

func walkMailbox(c *imap.Client, b string, l int) error {
	cmd, err := imap.Wait(c.List(b, "%"))
	if err != nil {
		return err
	}

	for pos, rsp := range cmd.Data {
		box := boxchar(pos, l, len(cmd.Data))
		fmt.Println(box, filepath.Base(rsp.MailboxInfo().Name))
		if rsp.MailboxInfo().Attrs["\\Haschildren"] {
			err = walkMailbox(c, rsp.MailboxInfo().Name+rsp.MailboxInfo().Delim, l+1)
			if err != nil {
				log.Printf("[ERR] While walking Mailboxes: %s\n", err)
				return err
			}
		}
	}
	return err
}

func boxchar(p, l, b int) string {
	var drawthis string
	switch {
	case p == b || p == 0 && l > 0:
		drawthis = "└─"
	case p == 0 && p < b:
		drawthis = "┌─"
	case p > 0 && p < b:
		drawthis = "├─"
	}
	if l > 0 {
		drawthis = "│" + strings.Repeat(" ", l) + drawthis
	}
	return drawthis
}

func main() {
	// imap.DefaultLogMask = imap.LogConn | imap.LogRaw
	fileconf := flag.String("conf", "", "Configuration file")
	list := flag.Bool("list", false, "List all mailboxes and exit")
	flag.Parse()
	raw, err := ioutil.ReadFile(*fileconf)
	if err != nil {
		log.Fatalf("[ERR] Can't read file: %s", err)
	}
	var conf NotifyConfig
	_ = json.Unmarshal(raw, &conf)

	if *list {
		client := newClient(conf)
		defer client.Logout(30 * time.Second)
		_ = walkMailbox(client, "", 0)
	} else {
		events := make(chan IDLEEvent, 100)
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		guard := guardian{
			mx: &sync.Mutex{},
			wg: &sync.WaitGroup{},
		}

		NewWatchMailBox(conf, events, quit, &guard)

		// Process incoming events from the mailboxes
		for rsp := range events {
			log.Printf("[DBG] Event %s for %s", rsp.EventType, rsp.Mailbox)
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
	}
}
