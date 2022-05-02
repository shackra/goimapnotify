package models

type UpdateKind int

const (
	FinishedWithFailure UpdateKind = iota
	FinishedSuccesfully
	NewMail
	DeletedMail
)

type Event struct {
	Kind    UpdateKind
	Mailbox string
	Error   error
}
