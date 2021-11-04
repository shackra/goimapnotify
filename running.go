package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	tag = "mail sync job"
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

	if rsp.OnNewMail == "SKIP" || rsp.OnNewMail == "" {
		logrus.Infof("[%s:%s] No scripts to be executed. Skipping..",
			rsp.Alias,
			rsp.Mailbox)
		return
	}
	key := rsp.Alias + rsp.Mailbox
	wait := time.Duration(r.wait) * time.Second
	format := fmt.Sprintf("[%s:%s] %%s syncing for %s (%s in the future)",
		rsp.Alias,
		rsp.Mailbox,
		time.Now().Add(wait).Format(time.RFC850),
		wait)

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
		logrus.Infof(format, "Scheduled")
		select {
		case <-r.timer[key].C:
			r.run(rsp)
		case <-done:
			//just get out
		}
	} else {
		logrus.Infof(format, "Rescheduled")
	}
}

func (r *RunningBox) run(rsp IDLEEvent) {
	if r.debug {
		logrus.Infof("[%s:%s] Running synchronization...",
			rsp.Alias,
			rsp.Mailbox)
	}

	if rsp.OnNewMail == "SKIP" || rsp.OnNewMail == "" {
		return
	}
	newmail := PrepareCommand(rsp.OnNewMail, rsp, r.debug)
	err := newmail.Run()
	if err != nil {
		logrus.Errorf("[%s:%s] OnNewMail command failed: %s",
			rsp.Alias,
			rsp.Mailbox,
			err)
	} else {
		if rsp.OnNewMailPost == "SKIP" ||
			rsp.OnNewMailPost == "" {
			return
		}
		newmailpost := PrepareCommand(rsp.OnNewMailPost, rsp, r.debug)
		err = newmailpost.Run()
		if err != nil {
			logrus.Errorf("[%s:%s] OnNewMailPost command failed: %s",
				rsp.Alias,
				rsp.Mailbox,
				err)
		}
	}
}
