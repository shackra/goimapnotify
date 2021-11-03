package main

// This file is part of goimapnotify
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
	"sync"

	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
)

// IDLEEvent models an IDLE event
type IDLEEvent = Box

// BoxEvent helps in communication between the box watch launcher and the box
// watching goroutines
type BoxEvent struct {
	Conf    NotifyConfig
	Mailbox Box
}

// WatchMailBox Keeps track of the IDLE state of one Mailbox
type WatchMailBox struct {
	client    *IMAPIDLEClient
	conf      NotifyConfig
	box       Box
	idleEvent chan<- IDLEEvent
	boxEvent  chan<- BoxEvent
	done      <-chan struct{}
}

func (w *WatchMailBox) Watch() {
	updates := make(chan client.Update)
	done := make(chan error, 1)

	if _, err := w.client.Select(w.box.Mailbox, true); err != nil {
		logrus.Fatalf("[%s:%s] Cannot select mailbox: %s",
		              w.box.Alias,
		              w.box.Mailbox,
		              err)
	}
	w.client.Updates = updates

	go func() {
		logrus.Infof("[%s:%s] Watching mailbox", w.box.Alias, w.box.Mailbox)
		done <- w.client.IdleWithFallback(nil, 0) // 0 = good default
		_ = w.client.Logout()
	}()

	// Block and process IDLE events
	for {
		select {
		case update := <-updates:
			m, ok := update.(*client.MailboxUpdate)
			if ok && m.Mailbox.Messages > 0 {
				// dispatch IDLE event to the main loop
				w.idleEvent <- w.box
			}
			// message deleted
			_, ok = update.(*client.ExpungeUpdate)
			if ok {
				w.idleEvent <- w.box
			}
		case <-w.done:
			// the main event loop is asking us to stop
			logrus.Warnf("[%s:%s] Stopping client watching mailbox",
			             w.box.Alias,
			             w.box.Mailbox)
			return
		case finished := <-done:
			logrus.Warnf("[%s:%s] Done watching mailbox",
			             w.box.Alias,
			             w.box.Mailbox)
			if finished != nil {
				w.boxEvent <- BoxEvent{Conf: w.conf, Mailbox: w.box}
			}
			return
		}
	}
}

// NewWatchBox creates a new instance of WatchMailBox and launch it
func NewWatchBox(c *IMAPIDLEClient, f NotifyConfig, m Box, i chan<- IDLEEvent,
                 b chan<- BoxEvent, q <-chan struct{}, wg *sync.WaitGroup) {
	w := WatchMailBox{
		client:    c,
		conf:      f,
		box:       m,
		idleEvent: i,
		boxEvent:  b,
		done:      q,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		w.Watch()
	}()
}
