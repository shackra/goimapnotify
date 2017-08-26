package main

// This file is part of goimapnotify
// Copyright (C) 2017  Jorge Javier Araya Navarro

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
	"log"

	"github.com/mxk/go-imap/imap"
)

func newClient(conf NotifyConfig) *imap.Client {
	var c *imap.Client
	var err error
	if conf.Port != 0 {
		c, err = imap.Dial(fmt.Sprintf("%s:%d", conf.Host, conf.Port))
	} else if conf.Port == 993 {
		c, err = imap.DialTLS(fmt.Sprintf("%s:%d", conf.Host, conf.Port), &tls.Config{
			ServerName:         conf.Host,
			InsecureSkipVerify: !conf.TLSOptions.RejectUnauthorized,
		})
	} else {
		c, err = imap.Dial(conf.Host)
	}

	if err != nil {
		log.Fatalf("[ERR] Cannot connect to %s:%d: %s", conf.Host, conf.Port, err)
	}

	// Enable encryption, if supported by the server
	if c.Caps["STARTTLS"] && conf.Port != 993 {
		// #nosec
		c.StartTLS(&tls.Config{
			ServerName:         conf.Host,
			InsecureSkipVerify: !conf.TLSOptions.RejectUnauthorized,
		})
	}

	// Authenticate
	if c.State() == imap.Login {
		_, err = c.Login(conf.Username, conf.Password)
	}

	if err != nil {
		log.Fatalf("[ERR] Can't login to %s with %s: %s", conf.Host, conf.Username, err)
	}
	log.Printf("Connected to %s\n", conf.Host)
	return c
}
