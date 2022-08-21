package main

// Execute scripts on events using IDLE imap command (Go version)
// Copyright (C) 2017-2021  Jorge Javier Araya Navarro

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
	"gitlab.com/shackra/goimapnotify/internal/providers/imap"
)

func main() {
	_, err := imap.New(&imap.Config{
		Host:     "127.0.0.1",
		Port:     1143,
		Username: "jorge@esavara.cr",
		Mailbox:  "INBOX",
		Opts: []imap.LoginOption{
			imap.WithDebug(true),
			imap.WithPassword("7SF5dZ_HpLIhLROcgIkmTQ"),
		},
	}, nil)
	if err != nil {
		panic(err)
	}
}
