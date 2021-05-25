package main

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type RunningBox struct {
	boxes []string
	m     *sync.RWMutex
	debug bool
}

func (r *RunningBox) RunOrIgnore(nm, nmp string, rsp IDLEEvent) {
	result := r.isBoxRunning(rsp)
	if result == -1 {
		r.m.Lock()
		defer r.m.Unlock()
		if r.debug {
			logrus.Infof("locking mailbox %s", rsp.Mailbox)
		}
		r.boxes = append(r.boxes, rsp.Mailbox)

		// execute commands for the Mailbox with an update
		go r.run(nm, nmp, rsp)
	} else if r.debug {
		logrus.Warnf("ignoring executing onNewMail & onNewMailPost commands for mailbox %s", rsp.Mailbox)
	}
}

func (r *RunningBox) run(nm, nmp string, rsp IDLEEvent) {
	// remove the Mailbox name from the list after completing the execution of
	// the system commands
	defer r.freeBox(rsp)

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

func (r *RunningBox) isBoxRunning(rsp IDLEEvent) int {
	for index, box := range r.boxes {
		if box == rsp.Mailbox {
			return index
		}
	}
	return -1
}

func (r *RunningBox) freeBox(rsp IDLEEvent) {
	index := r.isBoxRunning(rsp)
	r.m.Lock()
	defer r.m.Unlock()

	if r.debug {
		logrus.Infof("releasing mailbox %s", rsp.Mailbox)
	}

	if index > -1 {
		r.boxes[len(r.boxes)-1], r.boxes[index] = r.boxes[index], r.boxes[len(r.boxes)-1]
		r.boxes = r.boxes[:len(r.boxes)-1]
	}
}
