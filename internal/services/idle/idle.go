package idle

import "gitlab.com/shackra/goimapnotify/internal/services/models"

type CommanderEmailReceived interface {
	WhenReceived(mailbox string) error
	WhenReceivedPost(mailbox string) error
}

type CommanderEmailDeleted interface {
	WhenDeleted(mailbox string) error
	WhenDeletedPost(mailbox string) error
}

type IdleService struct {
	providers map[string]*models.IdleWatcher
	stop      chan struct{}
	events    chan models.Event
}

// Replace replaces an idleWatcher that suddenly stop running
func (i *IdleService) Replace(watcher *models.IdleWatcher) {
	name := watcher.GetMailbox()

	i.providers[name] = watcher

	go func() {
		watcher.WatchIdle(i.stop)
	}()
}

func (i *IdleService) Watch(receivedEmail CommanderEmailReceived, deletedEmail CommanderEmailDeleted) {
	i.start()

	select {
	case event := <-i.events:
		switch event.Kind {
		case models.ReceivedEmail:
			if err := receivedEmail.WhenReceived(event.Mailbox); err != nil {
				// TODO: log error
				break
			}
			if err := receivedEmail.WhenReceivedPost(event.Mailbox); err != nil {
				// TODO: log error
			}
		case models.DeletedEmail:
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

func (i *IdleService) start() {
	for name := range i.providers {
		provider := i.providers[name]
		go func() {
			provider.WatchIdle(i.stop)
		}()
	}
}

func New(watchers []*models.IdleWatcher, stop chan struct{}, events chan models.Event) *IdleService {
	providers := make(map[string]*models.IdleWatcher)

	for _, provider := range watchers {
		name := provider.GetMailbox()
		providers[name] = provider
	}

	return &IdleService{
		stop:      stop,
		providers: providers,
		events:    events,
	}
}
