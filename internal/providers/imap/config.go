package imap

type Config struct {
	Host     string
	Port     int
	Username string

	Mailbox string

	Opts []LoginOption
}
