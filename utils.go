package main

// This file is part of goimapnotify
// Copyright (C) 2017-2019  Jorge Javier Araya Navarro

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
	"fmt"
	"log"
	"path/filepath"
	"strings"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
)

func walkMailbox(c *client.Client, b string, l int) error {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List(b, "*", mailboxes)
	}()

	pos := 0
	for m := range mailboxes {
		box := boxchar(pos, l, len(m.Name))
		fmt.Println(box, filepath.Base(m.Name))
		// Check if mailbox has children mailboxes
		var godown bool
		for _, attr := range m.Attributes {
			if attr == "\\Haschildren" {
				godown = true
				break
			}
		}
		// if so, go down
		if godown {
			bErr := walkMailbox(c, m.Name+m.Delimiter, l+1)
			if bErr != nil {
				logrus.Errorf("cannot keep walking mailboxes: %s\n", bErr)
				return bErr
			}
		}
		pos += 1
	}
	if err := <-done; err != nil {
		return err
	}
	return nil
}

func boxchar(p, l, b int) string {
	var drawthis string
	switch {
	case p == b || p == 0 && l > 0:
		drawthis = "└─"
	case p == 0 && p < b:
		drawthis = "┌─"
	case p > 0 && p < b:
		drawthis = "├─"
	}
	if l > 0 {
		drawthis = "│" + strings.Repeat(" ", l) + drawthis
	}
	return drawthis
}

func retrievePasswordCmd(conf NotifyConfig) NotifyConfig {
	if conf.PasswordCMD != "" {
		cmd := PrepareCommand(conf.PasswordCMD, IDLEEvent{})
		// Avoid leaking the password
		cmd.Stdout = nil
		buf, err := cmd.Output()
		if err == nil {
			conf.Password = strings.Trim(string(buf[:]), "\n")
		} else {
			log.Fatalf("Can't retrieve password from command: %s", err)
		}
	}
	return conf
}

func retrieveUsernameCmd(conf NotifyConfig) NotifyConfig {
	if conf.PasswordCMD != "" {
		cmd := PrepareCommand(conf.UsernameCMD, IDLEEvent{})
		// Avoid leaking the username
		cmd.Stdout = nil
		buf, err := cmd.Output()
		if err == nil {
			conf.Username = strings.Trim(string(buf[:]), "\n")
		} else {
			log.Fatalf("Can't retrieve username from command: %s", err)
		}
	}
	return conf
}

func retrieveHostCmd(conf NotifyConfig) NotifyConfig {
	if conf.HostCMD != "" {
		cmd := PrepareCommand(conf.HostCMD, IDLEEvent{})
		// Avoid leaking the hostname
		cmd.Stdout = nil
		buf, err := cmd.Output()
		if err == nil {
			conf.Host = strings.Trim(string(buf[:]), "\n")
		} else {
			log.Fatalf("Can't retrieve host from command: %s", err)
		}
	}
	return conf
}
