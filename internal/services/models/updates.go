package models

type UpdateKind int

const (
	FinishedWithFailure UpdateKind = iota
	FinishedSuccesfully
	ArrivedEmail
	DeletedEmail
)

type Event struct {
	Kind    UpdateKind
	Mailbox string
	Error   error
}
