package main

// This file is part of goimapnotify
// Copyright (C) 2017-2023	Jorge Javier Araya Navarro

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
	"os/exec"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type RunningBox struct {
	debug bool
	wait  int
	/*
	 * Use map to create a different timer for each
	 * username-mailbox combination
	 */
	timer map[string]*time.Timer
	mutex map[string]*sync.RWMutex
}

func NewRunningBox(debug bool, wait int) *RunningBox {
	return &RunningBox{
		debug: debug,
		wait:  wait,
		timer: make(map[string]*time.Timer),
		mutex: make(map[string]*sync.RWMutex),
	}
}

func (r *RunningBox) schedule(rsp IDLEEvent, done <-chan struct{}) {
	l := logrus.WithField("alias", rsp.Alias).WithField("mailbox", rsp.Mailbox)
	if shouldSkip(rsp) {
		l.Warnf("No command for %q, skipping scheduling...", rsp.Reason)
		return
	}

	key := rsp.Alias + rsp.Mailbox
	wait := time.Duration(r.wait) * time.Second
	when := time.Now().Add(wait).Format(time.RFC850)
	format := fmt.Sprintf("%%s syncing %q for %s (%s in the future)", rsp.Reason, when, wait)

	r.mutex[key].Lock()
	_, exists := r.timer[key]
	main := true // main is true for the goroutine that will run sync
	if exists {
		// Stop should be called before Reset according to go docs
		if r.timer[key].Stop() {
			main = false // stopped running timer -> main is another goroutine
		}
		r.timer[key].Reset(wait)
	} else {
		r.timer[key] = time.NewTimer(wait)
	}
	r.mutex[key].Unlock()

	if main {
		l.Infof(format, "scheduled")
		select {
		case <-r.timer[key].C:
			r.run(rsp)
		case <-done:
			// just get out
		}
	} else {
		l.Infof(format, "rescheduled")
	}
}

func (r *RunningBox) run(rsp IDLEEvent) {
	l := logrus.WithField("alias", rsp.Alias).WithField("mailbox", rsp.Mailbox)
	if r.debug {
		l.Infoln("Running synchronization...")
	}

	var err error
	if rsp.Reason == NEWMAIL {
		err = prepareAndRun(rsp.OnNewMail, rsp.OnNewMailPost, rsp.Reason, rsp, r.debug)
	} else if rsp.Reason == DELETEDMAIL {
		err = prepareAndRun(rsp.OnDeletedMail, rsp.OnDeletedMailPost, rsp.Reason, rsp, r.debug)
	}

	if err != nil {
		logrus.WithError(err).Errorf("an error was encountered while executing commands for %q", rsp.Reason)
	}
}

func prepareAndRun(on, onpost string, kind EventType, event IDLEEvent, debug bool) error {
	callKind := "New"
	if kind == DELETEDMAIL {
		callKind = "Deleted"
	}

	if on == "SKIP" || on == "" {
		return nil
	}
	call := PrepareCommand(on, event)
	out, err := call.Output()
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok {
			logrus.Errorf("stderror: %q", string(exiterr.Stderr))
		}
		return fmt.Errorf("On%sMail command failed: %v", callKind, err)
	}
	logrus.Infof("stdout: %q", string(out))

	if onpost == "SKIP" || onpost == "" {
		return nil
	}

	call = PrepareCommand(onpost, event)
	out, err = call.Output()
	if err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if ok {
			logrus.Errorf("stderror: %q", string(exiterr.Stderr))
		}
		return fmt.Errorf("On%sMailPost command failed: %v", callKind, err)
	}
	logrus.Infof("stdout: %q", string(out))

	return nil
}

func shouldSkip(event IDLEEvent) bool {
	switch event.Reason {
	case NEWMAIL:
		return event.OnNewMail == "" || event.OnNewMail == "SKIP"
	case DELETEDMAIL:
		return event.OnDeletedMail == "" || event.OnDeletedMail == "SKIP"
	}

	return false
}
