package main

// This file is part of goimapnotify
// Copyright (C) 2017-2019  Jorge Javier Araya Navarro

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

type Event int

const (
	FirstStart Event = iota
	Bye
)

// IDLEEvent models an IDLE event
type IDLEEvent struct {
	Mailbox string
}

// BoxEvent helps in communication between the box watch launcher and the box
// watching goroutines
type BoxEvent struct {
	Type    Event
	Mailbox string
}

// WatchMailBox Keeps track of the IDLE state of one Mailbox
type WatchMailBox struct {
	client    *IMAPIDLEClient
	mailbox   string
	idleEvent chan<- IDLEEvent
	boxEvent  chan<- BoxEvent
	done      <-chan struct{}
}

func (w *WatchMailBox) Watch() {
	updates := make(chan client.Update)
	done := make(chan error, 1)

	if _, err := w.client.Select(w.mailbox, true); err != nil {
		logrus.Fatalf("cannot select mailbox %s, reason: %s", w.mailbox, err)
	}
	w.client.Updates = updates

	go func() {
		logrus.Infof("Watching mailbox %s", w.mailbox)
		done <- w.client.IdleWithFallback(nil, 0) // 0 = good default
		_ = w.client.Logout()
	}()

	// Block and process IDLE events
	for {
		select {
		case update := <-updates:
			_, ok := update.(*client.MailboxUpdate)
			if ok {
				// dispatch IDLE event to the main loop
				w.idleEvent <- IDLEEvent{Mailbox: w.mailbox}
			}
		case <-w.done:
			// the main event loop is asking us to stop
			logrus.Warn("Stopping client watching mailbox " + w.mailbox)
			return
		case finished := <-done:
			logrus.Warnf("Done watching mailbox %s", w.mailbox)
			if finished != nil {
				w.boxEvent <- BoxEvent{Type: Bye, Mailbox: w.mailbox}
			}
			return
		}
	}
}

// NewWatchBox creates a new instance of WatchMailBox and launch it
func NewWatchBox(c *IMAPIDLEClient, m string, i chan<- IDLEEvent, b chan<- BoxEvent, q <-chan struct{}, wg *sync.WaitGroup) {
	w := WatchMailBox{
		client:    c,
		mailbox:   m,
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
