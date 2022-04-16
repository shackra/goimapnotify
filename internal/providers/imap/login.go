package imap

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

const (
	// FIXME: make this configurable
	maxAttempts = 5
)

var (
	NoBearerSupport  = errors.New("server does not support OAuth Bearer token authentication")
	NoXOAuth2Support = errors.New("server does not support XOAuth2 authentication")
)

type CheckOAuthImpossible struct {
	cause error
}

func (c *CheckOAuthImpossible) Error() string {
	return fmt.Sprintf("is impossible to check support for OAuth: %v", c.cause)
}

func (c *CheckOAuthImpossible) Cause() error {
	return c.cause
}

func NewCheckOAuthImpossibleErr(err error) error {
	return &CheckOAuthImpossible{cause: err}
}

type loginOptions struct {
	debug              bool
	useOAuth           bool
	useTLS             bool
	insecureSkipVerify bool

	tokenCommand    string
	passwordCommand string
	usernameCommand string

	password string
}

type loginOption interface {
	apply(opts *loginOptions)
}

func dial(host string, port int, useTLS, skipVerify bool) (*client.Client, error) {
	var c *client.Client
	var err error = nil
	server := fmt.Sprintf("%s:%d", host, port)
	for attempt := 1; attempt < maxAttempts; attempt++ {
		if useTLS {
			c, err = client.DialTLS(server, &tls.Config{
				ServerName:         host,
				InsecureSkipVerify: skipVerify,
			})
		} else {
			c, err = client.Dial(server)
		}

		if err == nil {
			return c, nil
		} else {
			s := time.Duration(attempt) * 10 * time.Second
			// FIXME: display retry
			time.Sleep(s)
		}
	}

	return nil, err
}

func checkOAuthSupport(imapClient *client.Client) (bool, bool, error) {
	var (
		bearer  bool = false
		xoauth2 bool = false
	)

	if ok, err := imapClient.SupportAuth(sasl.OAuthBearer); err != nil {
		return bearer, xoauth2, NewCheckOAuthImpossibleErr(err)
	} else if ok {
		bearer = ok
	}

	if ok, err := imapClient.SupportAuth(Xoauth2); err != nil {
		return bearer, xoauth2, NewCheckOAuthImpossibleErr(err)
	} else if ok {
		xoauth2 = ok
	}

	return bearer, xoauth2, nil
}

func authWithBearer(imapClient *client.Client, username, password, host string, port int) (*client.Client, error) {
	sasl_oauth := &sasl.OAuthBearerOptions{
		Username: username,
		Token:    password,
		Host:     host,
		Port:     port,
	}
	sasl_client := sasl.NewOAuthBearerClient(sasl_oauth)
	err := imapClient.Authenticate(sasl_client)
	if err != nil {
		return nil, err
	}

	return imapClient, nil
}

func authWithXOAuth2(imapClient *client.Client, username, password string) (*client.Client, error) {
	sasl_xoauth2 := NewXoauth2Client(username, password)
	err := imapClient.Authenticate(sasl_xoauth2)
	if err != nil {
		return nil, err
	}

	return imapClient, nil
}

func authRegular(imapClient *client.Client, username, password string) (*client.Client, error) {
	err := imapClient.Login(username, password)

	if err != nil {
		return nil, err
	}

	return imapClient, nil
}

func newIMAPClient(host string, port int, username string, opts ...loginOption) (*client.Client, error) {
	// set default options
	options := loginOptions{
		debug:              false,
		useOAuth:           false,
		useTLS:             false,
		insecureSkipVerify: false,
		tokenCommand:       "",
		passwordCommand:    "",
		usernameCommand:    "",
		password:           "",
	}

	// apply any options passed
	for _, o := range opts {
		o.apply(&options)
	}

	imapClient, err := dial(host, port, options.useTLS, options.insecureSkipVerify)
	if err != nil {
		return nil, err
	}

	if options.debug {
		imapClient.SetDebug(os.Stdout)
	}

	if options.tokenCommand != "" {
		r, err := fetchToken(options.tokenCommand, options.debug)
		if err != nil {
			return nil, err
		}

		options.password = r
	}

	// NOTE: We don't overwrite the password field if OAuth is in use
	if options.passwordCommand != "" && !options.useOAuth {
		r, err := fetchPassword(options.passwordCommand, options.debug)
		if err != nil {
			return nil, err
		}

		options.password = r
	}

	if options.usernameCommand != "" {
		r, err := fetchUsername(options.usernameCommand, options.debug)
		if err != nil {
			return nil, err
		}

		username = r
	}

	if options.useOAuth {
		bearerSupported, xoauthSupported, err := checkOAuthSupport(imapClient)
		if err != nil {
			return nil, err
		}

		if !bearerSupported && !xoauthSupported {
			return nil, errors.New("XOAUTH2 and OAUTHBEARER are not supported by the server")
		}

		if bearerSupported {
			imapClient, err = authWithBearer(imapClient, username, options.password, host, port)
			if err != nil {
				return nil, err
			}
		} else if xoauthSupported {
			imapClient, err = authWithXOAuth2(imapClient, username, options.password)
			if err != nil {
				return nil, err
			}
		}
	} else {
		imapClient, err = authRegular(imapClient, username, options.password)
		if err != nil {
			return nil, err
		}
	}

	return imapClient, nil
}
