package imap

import (
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"github.com/pkg/errors"
)

type IdleClient struct {
	client  *client.Client
	idle    *idle.Client
	mailbox string
}

func (i *IdleClient) selectMailbox(mailbox string, updates chan client.Update) (chan error, error) {
	if _, err := i.client.Select(mailbox, true); err != nil {
		return nil, errors.Wrap(err, "idleClient SelectMailbox")
	}

	i.mailbox = mailbox
	i.client.Updates = updates

	finished := make(chan error)

	return finished, nil
}

func (i *IdleClient) WatchIdle(finished chan error, stop chan struct{}) error {
	// if finished sends non-nil error, something went wrong
	finished <- i.idle.IdleWithFallback(stop, 0) // 0 is a good default

	err := i.client.Logout()

	if err != nil {
		return errors.Wrapf(err, "[mailbox: %s] idleClient watchIdle client Logout", i.mailbox)
	}

	return nil
}

func newIdleClient(c *client.Client, mailbox string, updates chan client.Update) (*IdleClient, chan error, error) {
	id := &IdleClient{
		client:  c,
		idle:    idle.NewClient(c),
		mailbox: "",
	}

	finished, err := id.selectMailbox(mailbox, updates)
	if err != nil {
		return nil, nil, err
	}

	return id, finished, nil
}

func New(conf *Config, updates chan client.Update) (*IdleClient, chan error, error) {
	c, err := newLogin(conf.Host, conf.Port, conf.Username, conf.Opts...)
	if err != nil {
		return nil, nil, err
	}

	idleClient, finished, err := newIdleClient(c, conf.Mailbox, updates)
	if err != nil {
		return nil, nil, err
	}

	return idleClient, finished, nil
}
