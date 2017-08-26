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
	"reflect"
	"testing"
)

func TestBugArgs(t *testing.T) {
	var args = []string{"sh", "-c", "emacsclient -e '(something)'"}
	cmd := PrepareCommand("emacsclient -e '(something)'", IDLEEvent{})
	if !reflect.DeepEqual(cmd.Args, args) {
		t.Errorf("*cmd.Args are %+v, expected %+v", cmd.Args, args)
	}
}
