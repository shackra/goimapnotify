package main

// This file is part of goimapnotify
// Copyright (C) 2017-2022  Jorge Javier Araya Navarro

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
	"os"
	"time"

	idle "github.com/emersion/go-imap-idle"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-sasl"
	"github.com/sirupsen/logrus"
)

const (
	maxAttempts = 5
)

type IMAPIDLEClient struct {
	*client.Client
	*idle.IdleClient
}

func (i *IMAPIDLEClient) SetUpdates(u chan<- client.Update) {
	i.Updates = u
}

func newClient(conf *NotifyConfig) (c *client.Client, err error) {
	for attempt := 1; attempt < maxAttempts; attempt++ {
		if conf.TLS {
			c, err = client.DialTLS(fmt.Sprintf("%s:%d", conf.Host, conf.Port), &tls.Config{
				ServerName:         conf.Host,
				InsecureSkipVerify: !conf.TLSOptions.RejectUnauthorized,
			})
		} else {
			c, err = client.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
		}

		if err != nil {
			// wait a bit if something went wrong trying to connect to the IMAP server
			seconds := time.Duration(attempt) * 10 * time.Second
			logrus.WithError(err).WithField("host", conf.Host).WithField("port", conf.Port).Errorf("there was an error attempting to connect the email client, retrying in %s (attempt %d)", seconds, attempt)
			time.Sleep(seconds)
		} else {
			// things went okay, stop the loop
			break
		}
	}

	if err != nil {
		return c, err
	}

	// turn on debugging
	if conf.Debug {
		c.SetDebug(os.Stdout)
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
			err = c.Authenticate(sasl_client)
		} else if okXOAuth2 {
			sasl_xoauth2 := NewXoauth2Client(conf.Username, conf.Password)
			err = c.Authenticate(sasl_xoauth2)
		}
	} else {
		err = c.Login(conf.Username, conf.Password)
	}

	return c, err
}

func newIMAPIDLEClient(conf *NotifyConfig) (c *IMAPIDLEClient, err error) {
	confCMDExecuted, err := retrieveCmd(conf)
	if err != nil {
		return nil, err
	}

	i, err := newClient(confCMDExecuted)
	if err != nil {
		return c, err
	}
	return &IMAPIDLEClient{i, idle.NewClient(i)}, nil
}
