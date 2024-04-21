package main

// This file is part of goimapnotify
// Copyright (C) 2017-2023	Jorge Javier Araya Navarro

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

type EventType int

const (
	NEWMAIL EventType = iota + 1
	DELETEDMAIL
)

func (e EventType) String() string {
	switch e {
	case NEWMAIL:
		return "New Email"
	case DELETEDMAIL:
		return "Deleted Email"
	default:
		return "Unknown Event"
	}
}

// NotifyConfigLegacy holds the old configuration format
type NotifyConfigLegacy struct {
	Host              string           `json:"host"`
	HostCMD           string           `json:"hostCmd,omitempty"`
	Port              int              `json:"port"`
	TLS               bool             `json:"tls,omitempty"`
	TLSOptions        TLSOptionsStruct `json:"tlsOption"`
	Username          string           `json:"username"`
	UsernameCMD       string           `json:"usernameCmd,omitempty"`
	Password          string           `json:"password"`
	PasswordCMD       string           `json:"passwordCmd,omitempty"`
	XOAuth2           bool             `json:"xoauth2"`
	OnNewMail         string           `json:"onNewMail"`
	OnNewMailPost     string           `json:"onNewMailPost,omitempty"`
	OnDeletedMail     string           `json:"onDeletedMail,omitempty"`
	OnDeletedMailPost string           `json:"onDeletedMailPost,omitempty"`
	Debug             bool             `json:"-"`
	Boxes             []string         `json:"boxes"`
}

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host              string           `json:"host"`
	HostCMD           string           `json:"hostCmd,omitempty"`
	Port              int              `json:"port"`
	TLS               bool             `json:"tls,omitempty"`
	TLSOptions        TLSOptionsStruct `json:"tlsOption"`
	Username          string           `json:"username"`
	UsernameCMD       string           `json:"usernameCmd,omitempty"`
	Alias             string           `json:"alias"`
	Password          string           `json:"password"`
	PasswordCMD       string           `json:"passwordCmd,omitempty"`
	XOAuth2           bool             `json:"xoauth2"`
	OnNewMail         string           `json:"onNewMail"`
	OnNewMailPost     string           `json:"onNewMailPost,omitempty"`
	OnDeletedMail     string           `json:"onDeletedMail,omitempty"`
	OnDeletedMailPost string           `json:"onDeletedMailPost,omitempty"`
	Debug             bool             `json:"-"`
	Boxes             []Box            `json:"boxes"`
}

type TLSOptionsStruct struct {
	RejectUnauthorized bool `json:"reject_unauthorized"`
	STARTTLS           bool `json:"starttls"`
}

/*
Box stores all the necessary info needed to be passed in an
IDLEEvent handler routine, in order to schedule commands and
print informative messages
*/
type Box struct {
	Alias             string    `json:"-"`
	Mailbox           string    `json:"mailbox"`
	Reason            EventType `json:"-"`
	OnNewMail         string    `json:"OnNewMail"`
	OnNewMailPost     string    `json:"onNewMailPost"`
	OnDeletedMail     string    `json:"onDeletedMail"`
	OnDeletedMailPost string    `json:"onDeletedMailPost"`
	ExistingEmail     uint32    `json:"-"`
}

func legacyConverter(conf NotifyConfigLegacy) []NotifyConfig {
	var r []NotifyConfig
	var c NotifyConfig
	c.Host = conf.Host
	c.HostCMD = conf.HostCMD
	c.Port = conf.Port
	c.TLS = conf.TLS
	c.TLSOptions = conf.TLSOptions
	c.Username = conf.Username
	c.UsernameCMD = conf.UsernameCMD
	c.Password = conf.Password
	c.PasswordCMD = conf.PasswordCMD
	c.XOAuth2 = conf.XOAuth2
	c.OnNewMail = conf.OnNewMail
	c.OnNewMailPost = conf.OnNewMailPost
	c.OnDeletedMail = conf.OnDeletedMail
	c.OnDeletedMailPost = conf.OnDeletedMailPost
	for _, mailbox := range conf.Boxes {
		c.Boxes = append(c.Boxes, Box{Mailbox: mailbox})
	}
	return append(r, c)
}
