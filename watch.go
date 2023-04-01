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
	"time"

	"github.com/emersion/go-imap"
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

type idleClient interface {
	IdleWithFallback(<-chan struct{}, time.Duration) error
	Select(string, bool) (*imap.MailboxStatus, error)
	SetUpdates(chan<- client.Update)
	Logout() error
}

// WatchMailBox Keeps track of the IDLE state of one Mailbox
type WatchMailBox struct {
	client    idleClient
	conf      NotifyConfig
	box       Box
	idleEvent chan<- IDLEEvent
	boxEvent  chan<- BoxEvent
	done      <-chan struct{}
	l         *logrus.Entry
}

func (w *WatchMailBox) EmailArrived(m *client.MailboxUpdate) {
	if m.Mailbox.Messages > 0 {
		w.idleEvent <- w.box
	}
}

func (w *WatchMailBox) EmailDeleted(m *client.ExpungeUpdate) {
	w.idleEvent <- w.box
}

func (w *WatchMailBox) RestartWatchingBox(b BoxEvent) {
	w.l.Warn("restarting watch on mail box: %s", w.box.Mailbox)
	w.boxEvent <- b
}

func (w *WatchMailBox) Watch(whenExit func()) {
	updates := make(chan client.Update)
	done := make(chan error, 1)

	if _, err := w.client.Select(w.box.Mailbox, true); err != nil {
		w.l.WithError(err).Fatal("Cannot select mailbox")
	}
	w.client.SetUpdates(updates)

	go func() {
		w.l.Info("Watching mailbox")
		done <- w.client.IdleWithFallback(nil, 0) // 0 = good default
		_ = w.client.Logout()
	}()

	// called after this function exits
	defer whenExit()

	// Block and process IDLE events
	stop := false
	for !stop {
		select {
		case update := <-updates:
			switch event := update.(type) {
			case *client.MailboxUpdate:
				w.EmailArrived(event)
			case *client.ExpungeUpdate:
				w.EmailDeleted(event)
			}
		case <-w.done:
			// the main event loop is asking us to stop
			w.l.Info("Stopping client watching mailbox")
			stop = true
		case finished := <-done:
			w.l.Info("Done watching mailbox")
			if finished != nil {
				w.RestartWatchingBox(BoxEvent{
					Conf:    w.conf,
					Mailbox: w.box,
				})
			}
			stop = true
		}
	}
}

// NewWatchBox creates a new instance of WatchMailBox and launch it
func NewWatchBox(c idleClient, f NotifyConfig, m Box, i chan<- IDLEEvent,
	b chan<- BoxEvent, q <-chan struct{}, wg *sync.WaitGroup,
) {
	w := &WatchMailBox{
		client:    c,
		conf:      f,
		box:       m,
		idleEvent: i,
		boxEvent:  b,
		done:      q,
		l:         logrus.WithField("alias", m.Alias).WithField("mailbox", m.Mailbox),
	}

	wg.Add(1)
	go func() {
		w.Watch(func() {
			wg.Done()
		})
	}()
}
