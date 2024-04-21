package main

// This file is part of goimapnotify
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
	"strings"
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

	status, err := w.client.Select(w.box.Mailbox, true)
	if err != nil {
		if strings.Contains(err.Error(), "reason: Unknown Mailbox") {
			logrus.Warnf("[%s:%s] Cannot select mailbox: %v, skipped!", w.box.Alias, w.box.Mailbox, err)
			return
		}
		logrus.Fatalf("[%s:%s] Cannot select mailbox: %v", w.box.Alias, w.box.Mailbox, err)
	}
	w.box.ExistingEmail = status.Messages
	logrus.Debugf("[%s:%s] Existing mail: %d", w.box.Alias, w.box.Mailbox, w.box.ExistingEmail)

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
			if m, ok := update.(*client.MailboxUpdate); ok && m.Mailbox != nil {
				if m.Mailbox.Messages > w.box.ExistingEmail {
					// messages arrived
					w.idleEvent <- makeIDLEEvent(w.box, NEWMAIL)
				} else {
					// messages deleted
					w.idleEvent <- makeIDLEEvent(w.box, DELETEDMAIL)
				} // NOTE: What if the number is the same as before?
				logrus.Debugf("[%s:%s] Existing mail from %d to %d", w.box.Alias, w.box.Mailbox, w.box.ExistingEmail, m.Mailbox.Messages)
				w.box.ExistingEmail = m.Mailbox.Messages
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

func makeIDLEEvent(w IDLEEvent, e EventType) IDLEEvent {
	return IDLEEvent{
		Alias:             w.Alias,
		Mailbox:           w.Mailbox,
		Reason:            e,
		OnNewMail:         w.OnNewMail,
		OnNewMailPost:     w.OnNewMailPost,
		OnDeletedMail:     w.OnDeletedMail,
		OnDeletedMailPost: w.OnDeletedMailPost,
	}
}
