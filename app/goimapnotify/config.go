package goimapnotify

type Config struct {
	Host               string
	Port               int
	UseTLS             bool
	InsecureSkipVerify bool
	UseOAuth           bool

	Username  string
	Password  string
	Mailboxes []string

	TokenCommand    string
	PasswordCommand string
	UsernameCommand string

	Debug bool
}
