package idle

import "gitlab.com/shackra/goimapnotify/internal/services/models"

type idleWatcher interface {
	WatchIdle(chan struct{})
}

type commandNewEmail interface {
	NewEmail(mailbox string) error
	NewEmailPost(mailbox string) error
}

type commandDeletedEmail interface {
	DeletedEmail(mailbox string) error
	DeletedEmailPost(mailbox string) error
}

type idleService struct {
	providers []idleWatcher
	stop      chan struct{}
	events    chan models.Event
}

func (i *idleService) Watch(newEmail commandNewEmail, deletedEmail commandDeletedEmail) {
	// FIXME: better start these somewhere else, so that we can spin a new idleWatcher if needed
	for index := range i.providers {
		provider := i.providers[index]
		go func() {
			provider.WatchIdle(i.stop)
		}()
	}

	select {
	case event := <-i.events:
		switch event.Kind {
		case models.NewMail:
			if err := newEmail.NewEmail(event.Mailbox); err != nil {
				// TODO: log error
				break
			}
			if err := newEmail.NewEmailPost(event.Mailbox); err != nil {
				// TODO: log error
			}
		case models.DeletedMail:
			if err := deletedEmail.DeletedEmail(event.Mailbox); err != nil {
				// TODO: log error
				break
			}
			if err := deletedEmail.DeletedEmailPost(event.Mailbox); err != nil {
				// TODO: log error
			}
		}
	}
}

func New(watchers []idleWatcher, stop chan struct{}, events chan models.Event) *idleService {
	return &idleService{
		stop:      stop,
		providers: watchers,
		events:    events,
	}
}
