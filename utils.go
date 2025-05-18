package main

// This file is part of goimapnotify
// Copyright (C) 2017-2023  Jorge Javier Araya Navarro

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
	"strings"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/sirupsen/logrus"
)

func printDelimiter(c *client.Client) (int, error) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	i := 0
	m := <-mailboxes
	for range mailboxes {
		i += 1
	}
	if err := <-done; err != nil {
		return 0, err
	}

	fmt.Println("Hierarchy delimiter is:", m.Delimiter)
	return i, nil
}

func walkMailbox(c *client.Client, b string, l, max int) error {
	// FIXME: This can be done better
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List(b, "*", mailboxes)
	}()

	pos := 0
	for m := range mailboxes {
		box := boxchar(pos, l, max)
		fmt.Println(box, m.Name)
		pos += 1
		// Check if mailbox has children mailboxes
		for _, attr := range m.Attributes {
			if attr == "\\Haschildren" {
				err := walkMailbox(c, m.Name, l+1, max)
				if err != nil {
					logrus.WithError(err).Error("cannot keep walking mailboxes\n")
					return err
				}
				break
			}
		}
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
	case l > 0:
		drawthis = "│" + strings.Repeat(" ", l) + drawthis
	default:
		drawthis = "├─"
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
			logrus.WithError(err).Fatal("cannot retrieve password from command")
		}
	}
	return conf
}

func retrieveUsernameCmd(conf NotifyConfig) NotifyConfig {
	if conf.UsernameCMD != "" {
		cmd := PrepareCommand(conf.UsernameCMD, IDLEEvent{})
		// Avoid leaking the username
		cmd.Stdout = nil
		buf, err := cmd.Output()
		if err == nil {
			conf.Username = strings.Trim(string(buf[:]), "\n")
		} else {
			logrus.WithError(err).Fatal("cannot retrieve username from command")
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
			logrus.WithError(err).Fatal("cannot retrieve host from command")
		}
	}
	return conf
}

func retrieveCmd(conf NotifyConfig) NotifyConfig {
	if conf.PasswordCMD != "" {
		conf = retrievePasswordCmd(conf)
	}
	if conf.UsernameCMD != "" {
		conf = retrieveUsernameCmd(conf)
	}
	if conf.HostCMD != "" {
		conf = retrieveHostCmd(conf)
	}
	return conf
}

func setFromConfig(conf NotifyConfig, box Box) Box {
	if box.OnNewMail == "" {
		box.OnNewMail = conf.OnNewMail
	}
	err := compileTemplate(box.OnNewMail)
	if err != nil {
		logrus.WithError(err).Fatal("template is invalid for 'OnNewMail'")
	}
	if box.OnNewMailPost == "" {
		box.OnNewMailPost = conf.OnNewMailPost
	}
	err = compileTemplate(box.OnNewMailPost)
	if err != nil {
		logrus.WithError(err).Fatal("template is invalid for 'OnNewMailPost'")
	}

	// for deleted email
	if box.OnDeletedMail == "" {
		box.OnDeletedMail = conf.OnDeletedMail
	}
	err = compileTemplate(box.OnDeletedMail)
	if err != nil {
		logrus.WithError(err).Fatal("template is invalid for 'OnDeletedMail'")
	}
	if box.OnDeletedMailPost == "" {
		box.OnDeletedMailPost = conf.OnDeletedMailPost
	}
	err = compileTemplate(box.OnDeletedMailPost)
	if err != nil {
		logrus.WithError(err).Fatal("template is invalid for 'OnDeletedMailPost'")
	}

	box.Alias = conf.Alias
	return box
}
