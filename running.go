package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
)

var (
	tag = "mail sync job"
)

type RunningBox struct {
	wait       int
	scheduler  *gocron.Scheduler
	debug      bool
	currentJob *gocron.Job
	timeOfJon  time.Time
}

func NewRunningBox(wait int, debug bool) *RunningBox {
	r := &RunningBox{
		wait:      wait,
		scheduler: gocron.NewScheduler(time.UTC),
		debug:     debug,
	}

	// makes tags unique
	r.scheduler.TagsUnique()
	// start scheduler asynchronously
	r.scheduler.StartAsync()

	return r
}

func (r *RunningBox) RunOrIgnore(nm, nmp string, rsp IDLEEvent) {
	if r.currentJob != nil {
		// "restart" the job
		// This may not work as expected if the job was already started
		r.scheduler.Remove(r.currentJob)
		if time.Now().Before(r.timeOfJon) {
			logrus.Infof("scheduled job (hopefully) removed (expected to run at %s)", r.timeOfJon.Format(time.RFC850))
		}
	}
	var err error
	r.currentJob, err = r.scheduler.Every(r.wait).Seconds().LimitRunsTo(1).Tag(tag).SingletonMode().Do(func() {
		logrus.Infof("running mail synchronization scheduled for %s", r.timeOfJon.Format(time.RFC850))
		r.run(nm, nmp, rsp)
	})
	r.timeOfJon = time.Now().Add(time.Duration(r.wait) * time.Second)

	if err != nil {
		logrus.Errorf("error when scheduling mail sync: %s", err)
	} else {
		logrus.Infof("mail synchronization job schedule to run in %d second(s) or at %s", r.wait, r.timeOfJon.Format(time.RFC850))
	}
}

func (r *RunningBox) run(nm, nmp string, rsp IDLEEvent) {
	newmail := PrepareCommand(nm, rsp)
	dErr := newmail.Run()
	if dErr != nil {
		logrus.Errorf("OnNewMail command failed: %s", dErr)
	} else {
		newmailpost := PrepareCommand(nmp, rsp)
		pErr := newmailpost.Run()
		if pErr != nil {
			logrus.Errorf("OnNewMailPost command failed: %s", pErr)
		}
	}
}
