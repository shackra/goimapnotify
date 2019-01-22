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

	"github.com/mxk/go-imap/imap"
)

func walkMailbox(c *imap.Client, b string, l int) error {
	cmd, err := imap.Wait(c.List(b, "%"))
	if err != nil {
		return err
	}

	for pos, rsp := range cmd.Data {
		box := boxchar(pos, l, len(cmd.Data))
		fmt.Println(box, filepath.Base(rsp.MailboxInfo().Name))
		if rsp.MailboxInfo().Attrs["\\Haschildren"] {
			err = walkMailbox(c, rsp.MailboxInfo().Name+rsp.MailboxInfo().Delim, l+1)
			if err != nil {
				log.Printf("[ERR] While walking Mailboxes: %s\n", err)
				return err
			}
		}
	}
	return err
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
