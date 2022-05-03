package idle

import "gitlab.com/shackra/goimapnotify/internal/services/models"

type idleWatcher interface {
	GetMailbox() string
	WatchIdle(chan struct{})
}

type commanderEmailNew interface {
	WhenNew(mailbox string) error
	WhenNewPost(mailbox string) error
}

type commanderEmailDeleted interface {
	WhenDeleted(mailbox string) error
	WhenDeletedPost(mailbox string) error
}

type idleService struct {
	providers map[string]idleWatcher
	stop      chan struct{}
	events    chan models.Event
}

// Replace replaces an idleWatcher that suddenly stop running
func (i *idleService) Replace(watcher idleWatcher) {
	name := watcher.GetMailbox()

	i.providers[name] = watcher

	go func() {
		watcher.WatchIdle(i.stop)
	}()
}

func (i *idleService) Watch(newEmail commanderEmailNew, deletedEmail commanderEmailDeleted) {
	i.start()

	select {
	case event := <-i.events:
		switch event.Kind {
		case models.NewMail:
			if err := newEmail.WhenNew(event.Mailbox); err != nil {
				// TODO: log error
				break
			}
			if err := newEmail.WhenNewPost(event.Mailbox); err != nil {
				// TODO: log error
			}
		case models.DeletedMail:
			if err := deletedEmail.WhenDeleted(event.Mailbox); err != nil {
				// TODO: log error
				break
			}
			if err := deletedEmail.WhenDeletedPost(event.Mailbox); err != nil {
				// TODO: log error
			}
			// TODO: resend event in other cases to another channel
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
