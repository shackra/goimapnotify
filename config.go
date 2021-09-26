package main

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host       string `json:"host"`
	HostCMD    string `json:"hostCmd,omitempty"`
	Port       int    `json:"port"`
	TLS        bool   `json:"tls,omitempty"`
	TLSOptions struct {
		RejectUnauthorized bool `json:"reject_unauthorized"`
	} `json:"tlsOption"`
	Username      string   `json:"username"`
	UsernameCMD   string   `json:"usernameCmd,omitempty"`
	Password      string   `json:"password"`
	PasswordCMD   string   `json:"passwordCmd,omitempty"`
	XOAuth2       bool     `json:"xoauth2"`
	OnNewMail     string   `json:"onNewMail"`
	OnNewMailPost string   `json:"onNewMailPost,omitempty"`
	Wait          int      `json:"wait"`
	Debug         bool     `json:"-"`
	Boxes         []string `json:"boxes"`
}
