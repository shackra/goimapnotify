package imap

type debugOption bool

func (d debugOption) apply(opts *loginOptions) {
	opts.debug = bool(d)
}

func WithDebug(d bool) loginOption {
	return debugOption(d)
}

type useXOAuth2Option bool

func (u useXOAuth2Option) apply(opts *loginOptions) {
	opts.useOAuth = bool(u)
}

func WithXOAuth(x bool) loginOption {
	return useXOAuth2Option(x)
}

type tokenCommandOption string

func (t tokenCommandOption) apply(opts *loginOptions) {
	opts.tokenCommand = string(t)
}

func WithTokenCommand(c string) loginOption {
	return tokenCommandOption(c)
}

type passwordCommandOption string

func (t passwordCommandOption) apply(opts *loginOptions) {
	opts.passwordCommand = string(t)
}

func WithPasswordCommand(c string) loginOption {
	return passwordCommandOption(c)
}

type usernameCommandOption string

func (t usernameCommandOption) apply(opts *loginOptions) {
	opts.usernameCommand = string(t)
}

func WithUsernameCommand(c string) loginOption {
	return usernameCommandOption(c)
}

type passwordOption string

func (p passwordOption) apply(opts *loginOptions) {
	opts.password = string(p)
}

func WithPassword(p string) loginOption {
	return passwordOption(p)
}

type insecureSkipVerifyOption bool

func (i insecureSkipVerifyOption) apply(opts *loginOptions) {
	opts.insecureSkipVerify = bool(i)
}

func WithInsecureSkipVerify(i bool) loginOption {
	return insecureSkipVerifyOption(i)
}

type useTLSOption bool

func (i useTLSOption) apply(opts *loginOptions) {
	opts.useTLS = bool(i)
}

func WithTLS(i bool) loginOption {
	return useTLSOption(i)
}
