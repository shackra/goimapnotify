package main

// This file is part of goimapnotify
// Copyright (C) 2017-2024  Jorge Javier Araya Navarro

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"crypto/tls"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	imapid "github.com/emersion/go-imap-id"
	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
)

const (
	maxAttempts = 5
)

type IMAPIDLEClient struct {
	*client.Client
	*idle.IdleClient
}

func newClient(conf NotifyConfig) (c *client.Client, err error) {
	if conf.TLS && !conf.TLSOptions.STARTTLS {
		c, err = client.DialTLS(conf.Host+fmt.Sprintf(":%d", conf.Port), &tls.Config{
			ServerName:         conf.Host,
			InsecureSkipVerify: !conf.TLSOptions.RejectUnauthorized,
			MinVersion:         tls.VersionTLS13,
		})
	} else {
		c, err = client.Dial(conf.Host + fmt.Sprintf(":%d", conf.Port))
	}

	if err != nil {
		return c, fmt.Errorf("cannot dial to %s:%d, tls: %t, start TLS: %t. error: %w", conf.Host, conf.Port, conf.TLS, conf.TLSOptions.STARTTLS, err)
	}

	// turn on debugging
	if conf.Debug {
		pr, pw := io.Pipe()

		sigChan := make(chan os.Signal, 1)

		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			pr.Close() // close the pipe when the program is about to close
		}()

		go censorCredentials(pr, os.Stdout)

		c.SetDebug(pw)
	}

	if conf.TLS && conf.TLSOptions.STARTTLS {
		err = c.StartTLS(&tls.Config{
			ServerName:         conf.Host,
			InsecureSkipVerify: !conf.TLSOptions.RejectUnauthorized,
		})
		if err != nil {
			return nil, err
		}
	}

	idClient := imapid.NewClient(c)
	if _, err := idClient.ID(imapid.ID{
		imapid.FieldName:    "goimapnotify",
		imapid.FieldVersion: gittag,
	}); err != nil {
		return nil, err
	}

	if conf.XOAuth2 {
		okBearer, err := c.SupportAuth(sasl.OAuthBearer)
		if err != nil {
			return nil, CannotCheckSupportedAuthErr
		}
		okXOAuth2, err := c.SupportAuth(Xoauth2)
		if err != nil {
			return nil, CannotCheckSupportedAuthErr
		}

		if !okXOAuth2 && !okBearer {
			return nil, TokenAuthNotSupportedErr
		}

		if okBearer {
			sasl_oauth := &sasl.OAuthBearerOptions{
				Username: conf.Username,
				// Use something like https://github.com/google/oauth2l
				// in your passwordCmd to grab the token as a password
				Token: conf.Password,
				Host:  conf.Host,
				Port:  conf.Port,
			}
			sasl_client := sasl.NewOAuthBearerClient(sasl_oauth)
			if err := c.Authenticate(sasl_client); err != nil {
				return nil, err
			}
		} else if okXOAuth2 {
			sasl_xoauth2 := NewXoauth2Client(conf.Username, conf.Password)
			if err := c.Authenticate(sasl_xoauth2); err != nil {
				return nil, err
			}
		}
	} else {
		err = c.Login(conf.Username, conf.Password)
	}

	return c, err
}

func newIMAPIDLEClient(conf NotifyConfig) (c *IMAPIDLEClient, err error) {
	confCMDExecuted := retrieveCmd(conf)
	i, err := newClient(confCMDExecuted)
	if err != nil {
		return c, err
	}
	return &IMAPIDLEClient{i, idle.NewClient(i)}, nil
}
