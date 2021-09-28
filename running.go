package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	tag = "mail sync job"
)

type RunningBox struct {
	wait        int
	debug       bool
	ignoreCalls bool
	timer       *time.Timer
	timeOfJob   time.Time
}

func NewRunningBox(wait int, debug bool) *RunningBox {
	return &RunningBox{
		wait:  wait,
		debug: debug,
	}
}

func (r *RunningBox) RunOrIgnore(nm, nmp string, rsp IDLEEvent) {
	wait := time.Duration(r.wait) * time.Second
	newSchedule := time.Now().Add(wait)
	oldTime := r.timeOfJob

	messageReschedule := fmt.Sprintf("re-scheduling syncing from %s to %s (%s in the future)", oldTime.Format(time.RFC850), newSchedule.Format(time.RFC850), wait)
	messageScheduled := fmt.Sprintf("syncing scheduled for %s (%s from now)", newSchedule.Format(time.RFC850), wait)

	if r.timer == nil {
		r.timer = time.AfterFunc(wait, func() {
			r.run(nm, nmp, rsp)
		})
		logrus.Infoln(messageScheduled)
	} else if !r.ignoreCalls {
		// syncing is not running, reset timer

		// "For a Timer created with AfterFunc(d, f), Reset either reschedules
		// when f will run, in which case Reset returns true, or schedules f to
		// run again, in which case it returns false."
		rescheduled := r.timer.Reset(wait)
		r.timeOfJob = newSchedule
		if rescheduled {
			logrus.Infoln(messageReschedule)
		} else {
			logrus.Infoln(messageScheduled)
		}
	}

	if r.ignoreCalls {
		logrus.Warningf("Ignoring this request, scheduled job for %s is running", r.timeOfJob.Format(time.RFC850))
	}
}

func (r *RunningBox) run(nm, nmp string, rsp IDLEEvent) {
	r.ignoreCalls = true
	defer func() {
		// turn off the flag
		r.ignoreCalls = false
	}()

	if r.debug {
		logrus.Infoln("running sinchronization...")
	}

	newmail := PrepareCommand(nm, rsp, r.debug)
	newmailpost := PrepareCommand(nmp, rsp, r.debug)

	err := newmail.Run()
	if err != nil {
		logrus.Errorf("OnNewMail command failed: %s", err)
	} else {
		err = newmailpost.Run()
		if err != nil {
			logrus.Errorf("OnNewMailPost command failed: %s", err)
		}
	}
}
