package main

import (
	"errors"
	"fmt"
	"testing"

	imap "github.com/emersion/go-imap"
	imapClient "github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/subtle-byte/mockigo/match"
)

type WatchTestSuite struct {
	suite.Suite
	client idleClient
	box    Box
	l      *logrus.Entry
}

func (suite *WatchTestSuite) SetupSuite() {
	suite.l = logrus.WithField("testing", "")
}

func (suite *WatchTestSuite) SetupTest() {
	suite.box = Box{
		Alias:         "test",
		Mailbox:       "test",
		OnNewMail:     "",
		OnNewMailPost: "",
	}
}

func (suite *WatchTestSuite) TestDoneChannel() {
	client := NewIdleClientMock(suite.T())

	// setup client
	client.EXPECT().Select(match.Eq(suite.box.Mailbox), match.Eq(true)).Return(nil, nil)
	client.EXPECT().SetUpdates(match.Arg[chan<- imapClient.Update](match.Any[chan<- imapClient.Update]()))

	idleEvent := make(chan IDLEEvent)
	boxEvent := make(chan BoxEvent)
	done := make(chan struct{})
	doneInner := make(chan error)

	w := &WatchMailBox{
		client:    client,
		conf:      nil,
		box:       suite.box,
		idleEvent: idleEvent,
		boxEvent:  boxEvent,
		done:      done,
		l:         suite.l,
	}

	var result error = errors.New("unset")

	go func() {
		result = w.Watch(func() {}, doneInner, nil)
	}()

	// send the signal
	done <- struct{}{}

	assert.NoError(suite.T(), result, "unexpected error found")

	close(idleEvent)
	close(done)
	close(boxEvent)
	close(doneInner)
}

func (suite *WatchTestSuite) TestFinishedChannelRestart() {
	client := NewIdleClientMock(suite.T())

	// setup client
	client.EXPECT().Select(match.Eq(suite.box.Mailbox), match.Eq(true)).Return(nil, nil)
	client.EXPECT().SetUpdates(match.Arg[chan<- imapClient.Update](match.Any[chan<- imapClient.Update]()))

	idleEvent := make(chan IDLEEvent)
	boxEvent := make(chan BoxEvent)
	done := make(chan struct{})
	doneInner := make(chan error)

	w := &WatchMailBox{
		client:    client,
		conf:      nil,
		box:       suite.box,
		idleEvent: idleEvent,
		boxEvent:  boxEvent,
		done:      done,
		l:         suite.l,
	}

	go w.Watch(func() {}, doneInner, nil)

	doneInner <- errors.New("random error")

	event := <-boxEvent
	assert.Equal(suite.T(), event.Mailbox.Alias, suite.box.Alias)

	close(idleEvent)
	close(done)
	close(boxEvent)
	close(doneInner)
}

func (suite *WatchTestSuite) TestUpdatesChannelArrived() {
	client := NewIdleClientMock(suite.T())

	// setup client
	client.EXPECT().Select(match.Eq(suite.box.Mailbox), match.Eq(true)).Return(nil, nil)

	idleEvent := make(chan IDLEEvent)
	boxEvent := make(chan BoxEvent)
	done := make(chan struct{})
	doneInner := make(chan error)

	updatesChan := make(chan imapClient.Update)
	client.EXPECT().SetUpdates(match.Arg[chan<- imapClient.Update](match.Any[chan<- imapClient.Update]()))

	w := &WatchMailBox{
		client:    client,
		conf:      nil,
		box:       suite.box,
		idleEvent: idleEvent,
		boxEvent:  boxEvent,
		done:      done,
		l:         suite.l,
	}

	go w.Watch(func() {}, doneInner, updatesChan)
	// we send the event
	go func() {
		updatesChan <- &imapClient.MailboxUpdate{
			Mailbox: &imap.MailboxStatus{
				Name:     "INBOX",
				Messages: 1,
			},
		}
	}()

	event := <-idleEvent
	assert.Equal(suite.T(), "test", event.Alias)

	close(idleEvent)
	close(done)
	close(boxEvent)
	close(doneInner)
	close(updatesChan)
}

func (suite *WatchTestSuite) TestUpdatesChannelDeleted() {
	client := NewIdleClientMock(suite.T())

	// setup client
	client.EXPECT().Select(match.Eq(suite.box.Mailbox), match.Eq(true)).Return(nil, nil)

	idleEvent := make(chan IDLEEvent)
	boxEvent := make(chan BoxEvent)
	done := make(chan struct{})
	doneInner := make(chan error)

	updatesChan := make(chan imapClient.Update)
	client.EXPECT().SetUpdates(match.Arg[chan<- imapClient.Update](match.Any[chan<- imapClient.Update]()))

	w := &WatchMailBox{
		client:    client,
		conf:      nil,
		box:       suite.box,
		idleEvent: idleEvent,
		boxEvent:  boxEvent,
		done:      done,
		l:         suite.l,
	}

	go w.Watch(func() {}, doneInner, updatesChan)
	// we send the event
	go func() {
		updatesChan <- &imapClient.ExpungeUpdate{
			SeqNum: 0,
		}
	}()

	event := <-idleEvent
	assert.Equal(suite.T(), "test", event.Alias)

	close(idleEvent)
	close(done)
	close(boxEvent)
	close(doneInner)
	close(updatesChan)
}

func (suite *WatchTestSuite) TestWatchSelectFailed() {
	client := NewIdleClientMock(suite.T())

	// setup client
	client.EXPECT().Select(match.Eq(suite.box.Mailbox), match.Eq(true)).Return(nil, fmt.Errorf("something went wrong"))

	w := &WatchMailBox{
		client: client,
		conf:   nil,
		box:    suite.box,
		l:      suite.l,
	}

	err := w.Watch(func() {}, nil, nil)
	assert.Error(suite.T(), err)
}

func TestWatchTestSuite(t *testing.T) {
	suite.Run(t, new(WatchTestSuite))
}
