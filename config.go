package main

// This file is part of goimapnotify
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

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host       string `json:"host"`
	HostCMD    string `json:"hostCmd,omitempty"`
	Port       int    `json:"port"`
	TLS        bool   `json:"tls,omitempty"`
	TLSOptions struct {
		RejectUnauthorized bool `json:"reject_unauthorized"`
	} `json:"tlsOption"`
	Username      string   `json:"username"`
	UsernameCMD   string   `json:"usernameCmd,omitempty"`
	Password      string   `json:"password"`
	PasswordCMD   string   `json:"passwordCmd,omitempty"`
	XOAuth2       bool     `json:"xoauth2"`
	OnNewMail     string   `json:"onNewMail"`
	OnNewMailPost string   `json:"onNewMailPost,omitempty"`
	Wait          int      `json:"wait"`
	Debug         bool     `json:"-"`
	Boxes         []string `json:"boxes"`
}
