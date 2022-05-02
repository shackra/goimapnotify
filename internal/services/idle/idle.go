package idle

import "gitlab.com/shackra/goimapnotify/internal/services/models"

type idleWatcher interface {
	GetMailbox() string
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
	providers map[string]idleWatcher
	stop      chan struct{}
	events    chan models.Event
}

// Replace replaces a idleWatcher that suddenly stop running
func (i *idleService) Replace(watcher idleWatcher) {
	name := watcher.GetMailbox()

	i.providers[name] = watcher

	go func() {
		watcher.WatchIdle(i.stop)
	}()
}

func (i *idleService) Watch(newEmail commandNewEmail, deletedEmail commandDeletedEmail) {
	i.start()

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

func (i *idleService) start() {
	for name := range i.providers {
		provider := i.providers[name]
		go func() {
			provider.WatchIdle(i.stop)
		}()
	}
}

func New(watchers []idleWatcher, stop chan struct{}, events chan models.Event) *idleService {
	providers := make(map[string]idleWatcher)

	for _, provider := range watchers {
		name := provider.GetMailbox()
		providers[name] = provider
	}

	return &idleService{
		stop:      stop,
		providers: providers,
		events:    events,
	}
}
