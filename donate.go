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
	"fmt"
	"io"

	"github.com/fatih/color"
)

func printDonate(out io.Writer, padding int) {
	msg := donateMessage(padding)
	white := color.New(color.FgWhite, color.Bold)
	stars := white.Sprintf("%*s*****************************************************\n", padding, " ")

	fmt.Print(stars + msg + stars)
}

func donateMessage(padding int) string {
	msg := ""
	magenta := color.New(color.FgMagenta, color.Bold)
	white := color.New(color.FgWhite, color.Bold)

	// line
	msg += white.Sprintf("%*sIf you like this project, consider making a donation \n", padding, " ")
	// line
	msg += white.Sprintf("%*sto the author at ", padding, " ")
	msg += magenta.Sprint("https://ko-fi.com/K3K1XEZCQ")
	msg += white.Sprint(" :D\n")
	// end

	return msg
}
