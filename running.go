package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
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

	var err error
	if rsp.Reason == NEWMAIL {
		err = prepareAndRun(rsp.OnNewMail, rsp.OnNewMailPost, rsp.Reason, rsp, r.debug)
	} else if rsp.Reason == DELETEDMAIL {
		err = prepareAndRun(rsp.OnDeletedMail, rsp.OnDeletedMailPost, rsp.Reason, rsp, r.debug)
	}

	if err != nil {
		logrus.Error(err)
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
	call := PrepareCommand(on, event, debug)
	err := call.Run()

	if err != nil {
		return fmt.Errorf("[%s:%s] On%sMail command failed: %v", event.Alias, event.Mailbox, callKind, err)
	}

	if onpost == "SKIP" || onpost == "" {
		return nil
	}
	call = PrepareCommand(onpost, event, debug)
	err = call.Run()
	if err != nil {
		return fmt.Errorf("[%s:%s] On%sMailPost command failed: %v", event.Alias, event.Mailbox, callKind, err)
	}

	return nil
}
