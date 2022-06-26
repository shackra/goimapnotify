package imap

import (
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"github.com/pkg/errors"

	"gitlab.com/shackra/goimapnotify/internal/services/models"
)

type IdleClient struct {
	client   *client.Client
	idle     *idle.Client
	mailbox  string
	updates  chan client.Update
	finished chan error
	events   chan models.Event
}

func (i *IdleClient) selectMailbox(mailbox string) error {
	if _, err := i.client.Select(mailbox, true); err != nil {
		return errors.Wrap(err, "idleClient SelectMailbox")
	}

	i.mailbox = mailbox
	i.client.Updates = i.updates

	i.finished = make(chan error)

	return nil
}

func (i *IdleClient) GetMailbox() string {
	return i.mailbox
}

func (i *IdleClient) WatchIdle(stop chan struct{}) {
	go func() {
		// if finished sends non-nil error, something went wrong
		i.finished <- i.idle.IdleWithFallback(stop, 0) // 0 is a good default

		err := i.client.Logout()

		if err != nil {
			// TODO: log error here
		}
	}()

	select {
	case update := <-i.updates:
		if m, ok := update.(*client.MailboxUpdate); ok && m.Mailbox.Messages > 0 {
			i.events <- models.Event{
				Kind:    models.ReceivedEmail,
				Mailbox: i.mailbox,
			}
		}
		if _, ok := update.(*client.ExpungeUpdate); ok {
			i.events <- models.Event{
				Kind:    models.DeletedEmail,
				Mailbox: i.mailbox,
			}
		}
	case finish := <-i.finished:
		if finish != nil {
			i.events <- models.Event{
				Kind:    models.FinishedWithFailure,
				Error:   errors.Wrap(finish, "idleClient WatchIdle finish"),
				Mailbox: i.mailbox,
			}
		} else {
			i.events <- models.Event{
				Kind:    models.FinishedSuccesfully,
				Mailbox: i.mailbox,
			}
		}
	case <-stop:
		if i.updates != nil {
			close(i.updates)
		}
		return
	}
}

func newIdleClient(c *client.Client, mailbox string, updates chan client.Update, events chan models.Event) (*IdleClient, error) {
	id := &IdleClient{
		client:  c,
		idle:    idle.NewClient(c),
		mailbox: "",
		events:  events,
	}

	err := id.selectMailbox(mailbox)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func New(conf *Config, events chan models.Event) (*IdleClient, error) {
	c, err := newLogin(conf.Host, conf.Port, conf.Username, conf.Opts...)
	if err != nil {
		return nil, err
	}

	internalUpdates := make(chan client.Update)

	idleClient, err := newIdleClient(c, conf.Mailbox, internalUpdates, events)
	if err != nil {
		return nil, err
	}

	return idleClient, nil
}
