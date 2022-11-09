package commands

type ShellReceived struct{}

func (s *ShellReceived) WhenReceived(mailbox string) error {
	return nil
}

func (s *ShellReceived) WhenReceivedPost(mailbox string) error {
	return nil
}

type ShellDeleted struct{}

func (s *ShellDeleted) WhenDeleted(mailbox string) error {
	return nil
}

func (s *ShellDeleted) WhenDeletedPost(mailbox string) error {
	return nil
}

func New(received, receivedPost, deleted, deletedPost string) (*ShellReceived, *ShellDeleted) {
	return &ShellReceived{}, &ShellDeleted{}
}
