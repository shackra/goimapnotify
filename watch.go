package main

// This file is part of goimapnotify
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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mxk/go-imap/imap"
)

// WatchMailBox Keeps track of the IDLE state of one Mailbox
type WatchMailBox struct {
	c     *imap.Client
	cmd   *imap.Command
	rsp   *imap.Response
	event chan<- IDLEEvent
	quit  <-chan os.Signal
}

// NewWatchMailBox create a list of WatchMailBoxes and start them in parallel
func NewWatchMailBox(conf NotifyConfig, event chan IDLEEvent, quit chan os.Signal, g *guardian) {
	for _, box := range conf.Boxes {
		var err error
		var watch WatchMailBox
		mailbox := box
		if conf.Port != 0 {
			watch.c, err = imap.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
		} else if conf.Port == 993 {
			watch.c, err = imap.DialTLS(fmt.Sprintf("%s:%d", conf.Host, conf.Port), nil)
		} else {
			watch.c, err = imap.Dial(conf.Host)
		}

		if err != nil {
			log.Fatalf("[ERR] Cannot connect to %s:%d: %s", conf.Host, conf.Port, err)
		}

		// Enable encryption, if supported by the server
		if watch.c.Caps["STARTTLS"] && conf.TLS && conf.Port != 993 {
			watch.c.StartTLS(nil)
		}

		// Authenticate
		if watch.c.State() == imap.Login {
			_, err = watch.c.Login(conf.Username, conf.Password)
		}

		if err != nil {
			log.Fatalf("[ERR] Can't login to %s with %s: %s", conf.Host, conf.Username, err)
		}

		// Include channels
		watch.quit = quit
		watch.event = event

		_, err = watch.c.Select(box, true)
		if err != nil {
			log.Fatalf("[ERR] Can't SELECT mailbox %s", box)
		}
		watch.c.Data = nil

		go func() {
			g.Add(1)
			defer g.Done()
			defer watch.c.Logout(30 * time.Second)
			idle, err := watch.c.Idle()
			if err != nil {
				log.Fatalf("[ERR] Can't start IDLE command: %s", err)
			}

			for idle.InProgress() {
				select {
				case <-watch.quit:
					idle, _ = watch.c.IdleTerm()
					log.Printf("[INF] Stopping watcher for box %s", mailbox)
				default:
					err = watch.c.Recv(1 * time.Second)
					// Process unilateral server data
					if err == nil {
						for _, watch.rsp = range watch.c.Data {
							// Create events and send them through
							// the channel
							var rsp = IDLEEvent{
								Mailbox:   mailbox,
								EventType: watch.rsp.Label,
							}
							watch.event <- rsp
						}
						watch.c.Data = nil
					}

				}
			}
			g.Close(watch.event)
		}()
	}
}
