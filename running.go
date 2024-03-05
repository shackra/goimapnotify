package main

import (
	"fmt"
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
	if rsp.OnNewMail == "SKIP" || rsp.OnNewMail == "" {
		l.Infoln("No scripts to be executed. Skipping...")
		return
	}
	key := rsp.Alias + rsp.Mailbox
	wait := time.Duration(r.wait) * time.Second
	format := fmt.Sprintf("%%s syncing for %s (%s in the future)",
		time.Now().Add(wait).Format(time.RFC850), wait)

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
		l.Infof(format, "Scheduled")
		select {
		case <-r.timer[key].C:
			r.run(rsp)
		case <-done:
			//just get out
		}
	} else {
		l.Infof(format, "Rescheduled")
	}
}

func (r *RunningBox) run(rsp IDLEEvent) {
	l := logrus.WithField("alias", rsp.Alias).WithField("mailbox", rsp.Mailbox)
	if r.debug {
		l.Infoln("Running synchronization...")
	}

	if rsp.OnNewMail == "SKIP" || rsp.OnNewMail == "" {
		return
	}
	newmail := PrepareCommand(rsp.OnNewMail, rsp, r.debug)
	err := newmail.Run()
	if err != nil {
		l.WithError(err).Errorln("OnNewMail command failed")
	} else {
		if rsp.OnNewMailPost == "SKIP" ||
			rsp.OnNewMailPost == "" {
			return
		}
		newmailpost := PrepareCommand(rsp.OnNewMailPost, rsp, r.debug)
		err = newmailpost.Run()
		if err != nil {
			l.WithError(err).Errorln("OnNewMailPost command failed")
		}
	}
}
