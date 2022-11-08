package goimapnotify

import (
	"gitlab.com/shackra/goimapnotify/internal/providers/imap"
	"gitlab.com/shackra/goimapnotify/internal/services/commands"
	"gitlab.com/shackra/goimapnotify/internal/services/idle"
	"gitlab.com/shackra/goimapnotify/internal/services/models"
)

type App struct {
	events   chan models.Event
	stop     chan struct{}
	service  *idle.IdleService
	received idle.CommanderEmailReceived
	deleted  idle.CommanderEmailDeleted
}

func (a *App) Start() {
	a.service.Watch(a.received, a.deleted)
}

func New(conf *Config) (*App, error) {
	received, deleted := commands.New(conf.ReceivedEmailCommand, conf.ReceivedEmailPostCommand, conf.DeletedEmailCommand, conf.DeletedEmailPostCommand)
	app := &App{
		events:   make(chan models.Event),
		stop:     make(chan struct{}),
		received: received,
		deleted:  deleted,
	}

	clients := make([]*models.IdleWatcher, len(conf.Mailboxes))

	for clientIndex := 0; clientIndex < len(conf.Mailboxes); clientIndex++ {
		idleConf := imap.Config{
			Host:     conf.Host,
			Port:     conf.Port,
			Username: conf.Username,
			Mailbox:  conf.Mailboxes[clientIndex],
			Opts: []imap.LoginOption{
				imap.WithDebug(conf.Debug),
				imap.WithTLS(conf.UseTLS),
				imap.WithInsecureSkipVerify(conf.InsecureSkipVerify),
				imap.WithPassword(conf.Password),
				imap.WithTokenCommand(conf.TokenCommand),
				imap.WithUsernameCommand(conf.UsernameCommand),
				imap.WithPasswordCommand(conf.PasswordCommand),
				imap.WithXOAuth(conf.UseOAuth),
			},
		}

		idleClient, err := imap.New(&idleConf, app.events)
		if err != nil {
			return nil, err
		}

		clients[clientIndex] = models.NewIdleWatcher(idleClient)
	}

	app.service = idle.New(clients, app.stop, app.events)

	return app, nil
}