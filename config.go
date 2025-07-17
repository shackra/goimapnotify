package main

import (
	"bytes"
	"fmt"
	"slices"
	"text/template"

	"github.com/emersion/go-imap"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// This file is part of goimapnotify
// Copyright (C) 2017-2025	Jorge Javier Araya Navarro

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
        FLAGCHANGED
)

func (e EventType) String() string {
	switch e {
	case NEWMAIL:
		return "New Email"
	case DELETEDMAIL:
		return "Deleted Email"
	case FLAGCHANGED:
		return "Changed Flag on Email"
	default:
		return "Unknown Event"
	}
}

// compileTemplate tests that the string template is valid, if any was
// provided.
func compileTemplate(i string) error {
	t, err := template.New("test").Parse(i)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)

	input := IDLEEvent{
		Alias:   "example@example.com",
		Mailbox: "Inbox",
	}

	return t.Execute(buf, input)
}

type Configuration struct {
	Configurations []NotifyConfig `json:"configurations" yaml:"configurations"`
}

// ConfigurationLegacy holds the old configuration format
type ConfigurationLegacy struct {
	Host              string           `yaml:"host"              json:"host"`
	HostCMD           string           `yaml:"hostCMD"           json:"hostCMD"`
	Port              int              `yaml:"port"              json:"port"`
	TLS               bool             `yaml:"tls"               json:"tls"`
	TLSOptions        TLSOptionsStruct `yaml:"tlsOptions"        json:"tlsOptions"`
	IDLELogoutTimeout int              `yaml:"idleLogoutTimeout" json:"idleLogoutTimeout"`
	EnableIDCommand   bool             `yaml:"enableIDCommand"   json:"enableIDCommand"`
	Username          string           `yaml:"username"          json:"username"`
	UsernameCMD       string           `yaml:"usernameCMD"       json:"usernameCMD"`
	Password          string           `yaml:"password"          json:"password"`
	PasswordCMD       string           `yaml:"passwordCMD"       json:"passwordCMD"`
	XOAuth2           bool             `yaml:"xoAuth2"           json:"xoAuth2"`
	OnNewMail         string           `yaml:"onNewMail"         json:"onNewMail"`
	OnNewMailPost     string           `yaml:"onNewMailPost"     json:"onNewMailPost"`
	OnChangedMail     string           `yaml:"onChangedMail"     json:"onChangedMail"`
	OnChangedMailPost string           `yaml:"onChangedMailPost" json:"onChangedMailPost"`
	OnDeletedMail     string           `yaml:"onDeletedMail"     json:"onDeletedMail"`
	OnDeletedMailPost string           `yaml:"onDeletedMailPost" json:"onDeletedMailPost"`
	Boxes             []string         `yaml:"boxes"             json:"boxes"`
}

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host              string           `yaml:"host"              json:"host"`
	HostCMD           string           `yaml:"hostCMD"           json:"hostCMD"`
	Port              int              `yaml:"port"              json:"port"`
	TLS               bool             `yaml:"tls"               json:"tls"`
	TLSOptions        TLSOptionsStruct `yaml:"tlsOptions"        json:"tlsOptions"`
	IDLELogoutTimeout int              `yaml:"idleLogoutTimeout" json:"idleLogoutTimeout"`
	EnableIDCommand   bool             `yaml:"enableIDCommand"   json:"enableIDCommand"`
	Username          string           `yaml:"username"          json:"username"`
	UsernameCMD       string           `yaml:"usernameCMD"       json:"usernameCMD"`
	Alias             string           `yaml:"alias"             json:"alias"`
	Password          string           `yaml:"password"          json:"password"`
	PasswordCMD       string           `yaml:"passwordCMD"       json:"passwordCMD"`
	XOAuth2           bool             `yaml:"xoAuth2"           json:"xoAuth2"`
	OnNewMail         string           `yaml:"onNewMail"         json:"onNewMail"`
	OnNewMailPost     string           `yaml:"onNewMailPost"     json:"onNewMailPost"`
	OnChangedMail     string           `yaml:"onChangedMail"     json:"onChangedMail"`
	OnChangedMailPost string           `yaml:"onChangedMailPost" json:"onChangedMailPost"`
	OnDeletedMail     string           `yaml:"onDeletedMail"     json:"onDeletedMail"`
	OnDeletedMailPost string           `yaml:"onDeletedMailPost" json:"onDeletedMailPost"`
	Boxes             []Box            `yaml:"boxes"             json:"boxes"`
}

type TLSOptionsStruct struct {
	RejectUnauthorized bool `yaml:"rejectUnauthorized" json:"rejectUnauthorized"`
	STARTTLS           bool `yaml:"starttls"           json:"starttls"`
}

/*
Box stores all the necessary info needed to be passed in an
IDLEEvent handler routine, in order to schedule commands and
print informative messages
*/
type Box struct {
	Alias             string    `json:"-"                 yaml:"-"`
	Mailbox           string    `json:"mailbox"           yaml:"mailbox"`
	Reason            EventType `json:"-"                 yaml:"-"`
	OnNewMail         string    `json:"onNewMail"         yaml:"onNewMail"`
	OnNewMailPost     string    `json:"onNewMailPost"     yaml:"onNewMailPost"`
	OnChangedMail     string    `json:"onChangedMail"     yaml:"onChangedMail"`
	OnChangedMailPost string    `json:"onChangedMailPost" yaml:"onChangedMailPost"`
	OnDeletedMail     string    `json:"onDeletedMail"     yaml:"onDeletedMail"`
	OnDeletedMailPost string    `json:"onDeletedMailPost" yaml:"onDeletedMailPost"`
	ExistingEmail     uint32    `json:"-"                 yaml:"-"`
}

func legacyConverter(conf ConfigurationLegacy) []NotifyConfig {
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
	c.OnChangedMail = conf.OnChangedMail
	c.OnChangedMailPost = conf.OnChangedMailPost
	c.OnDeletedMail = conf.OnDeletedMail
	c.OnDeletedMailPost = conf.OnDeletedMailPost
	c.IDLELogoutTimeout = conf.IDLELogoutTimeout
	c.EnableIDCommand = conf.EnableIDCommand
	for _, mailbox := range conf.Boxes {
		c.Boxes = append(c.Boxes, Box{Mailbox: mailbox})
	}
	return append(r, c)
}

func loadConfiguration(path string) (*Configuration, error) {
	var topConfiguration Configuration
	if err := viper.Unmarshal(&topConfiguration); err != nil {
		return nil, fmt.Errorf("can't parse the configuration: %q, error: %v", path, err)
	}

	if topConfiguration.Configurations == nil {
		var legacy ConfigurationLegacy
		if err := viper.UnmarshalExact(&legacy); err != nil {
			return nil, fmt.Errorf(
				"can't parse the configuration in 'legacy' format: %s, error: %v",
				path,
				err,
			)
		}

		logrus.Info("legacy format configuration detected")
		topConfiguration.Configurations = legacyConverter(legacy)
	}

	if len(topConfiguration.Configurations) > 0 &&
		(topConfiguration.Configurations[0].Host == "" && topConfiguration.Configurations[0].HostCMD == "") {
		return nil, fmt.Errorf(
			"configuration file %q is empty or have invalid configuration format",
			path,
		)
	}

	for account := range topConfiguration.Configurations {
		topConfiguration.Configurations[account] = retrieveCmd(
			topConfiguration.Configurations[account],
		)
		if topConfiguration.Configurations[account].Alias == "" {
			topConfiguration.Configurations[account].Alias = topConfiguration.Configurations[account].Username
		}
		if logrus.GetLevel() == logrus.DebugLevel {
			topConfiguration.Configurations[account].Alias = "<?>"
		}

		conf := topConfiguration.Configurations[account]

		// If there is no mailboxes, watch over all mailboxes of the account
		if len(conf.Boxes) == 0 {
			client, err := newIMAPIDLEClient(conf)
			if err != nil {
				return nil, fmt.Errorf(
					"account %q, failed to create IMAP client, error: %w",
					conf.Username,
					err,
				)
			}
			// nolint
			defer client.Logout()

			// NOTE(shackra): Having to do this is really disgusting, v2 offers a better way for listing mailboxes. I should consider updating.
			ch := make(chan *imap.MailboxInfo)
			go func() {
				err := client.List("", "*", ch)
				if err != nil {
					logrus.WithError(err).
						WithField("account", conf.Username).
						Fatal("failed to list all mailboxes")
				}
			}()

			for mailbox := range ch {
				// Ignore mailboxes with attributes `\All` and `\Noselect`
				if slices.Contains(mailbox.Attributes, "\\All") ||
					slices.Contains(mailbox.Attributes, "\\Noselect") {
					continue
				}

				box := setFromConfig(conf, Box{
					Mailbox: mailbox.Name,
				})
				topConfiguration.Configurations[account].Boxes = append(
					topConfiguration.Configurations[account].Boxes,
					box,
				)
			}
		} else {
			// replace all listed mailboxes with the same mailboxes carrying values from the configuration
			for mailbox := range topConfiguration.Configurations[account].Boxes {
				topConfiguration.Configurations[account].Boxes[mailbox] = setFromConfig(conf, topConfiguration.Configurations[account].Boxes[mailbox])
			}
		}
	}

	return &topConfiguration, nil
}
