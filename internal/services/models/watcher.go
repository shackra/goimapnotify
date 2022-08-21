package models

type idleWatcher interface {
	GetMailbox() string
	WatchIdle(chan struct{})
}

type IdleWatcher struct {
	watcher idleWatcher
}

func (i *IdleWatcher) GetMailbox() string {
	return i.watcher.GetMailbox()
}

func (i *IdleWatcher) WatchIdle(stop chan struct{}) {
	i.watcher.WatchIdle(stop)
}

func NewIdleWatcher(watcher idleWatcher) *IdleWatcher {
	return &IdleWatcher{watcher: watcher}
}
