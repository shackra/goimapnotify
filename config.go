package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

// This file is part of goimapnotify
// Copyright (C) 2017-2021	Jorge Javier Araya Navarro

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

// NotifyConfigLegacy holds the old configuration format
type NotifyConfigLegacy struct {
	Host          string           `json:"host"`
	HostCMD       string           `json:"hostCmd,omitempty"`
	Port          int              `json:"port"`
	TLS           bool             `json:"tls,omitempty"`
	TLSOptions    TLSOptionsStruct `json:"tlsOption"`
	Username      string           `json:"username"`
	UsernameCMD   string           `json:"usernameCmd,omitempty"`
	Password      string           `json:"password"`
	PasswordCMD   string           `json:"passwordCmd,omitempty"`
	XOAuth2       bool             `json:"xoauth2"`
	OnNewMail     string           `json:"onNewMail"`
	OnNewMailPost string           `json:"onNewMailPost,omitempty"`
	Debug         bool             `json:"-"`
	Boxes         []string         `json:"boxes"`
}

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host          string           `json:"host"`
	HostCMD       string           `json:"hostCmd,omitempty"`
	Port          int              `json:"port"`
	TLS           bool             `json:"tls,omitempty"`
	TLSOptions    TLSOptionsStruct `json:"tlsOption"`
	Username      string           `json:"username"`
	UsernameCMD   string           `json:"usernameCmd,omitempty"`
	Alias         string           `json:"alias"`
	Password      string           `json:"password"`
	PasswordCMD   string           `json:"passwordCmd,omitempty"`
	XOAuth2       bool             `json:"xoauth2"`
	OnNewMail     string           `json:"onNewMail"`
	OnNewMailPost string           `json:"onNewMailPost,omitempty"`
	Debug         bool             `json:"-"`
	Boxes         []Box            `json:"boxes"`
}

type TLSOptionsStruct struct {
	RejectUnauthorized bool `json:"reject_unauthorized"`
}

/*
Box stores all the necessary info needed to be passed in an
IDLEEvent handler routine, in order to schedule commands and
print informative messages
*/
type Box struct {
	Alias         string `json:"-"`
	Mailbox       string `json:"mailbox"`
	OnNewMail     string `json:"onNewMail"`
	OnNewMailPost string `json:"onNewMailPost"`
}

type configurationError struct {
	errors []error
}

func (c *configurationError) Error() string {
	var found []string
	for _, v := range c.errors {
		found = append(found, v.Error())
	}
	return fmt.Sprintf("When trying to load the configuration we found the following issues > %s", strings.Join(found, "; "))
}

func (c *configurationError) Push(err error) {
	c.errors = append(c.errors, err)
}

func newConfigurationError(err error) *configurationError {
	return &configurationError{[]error{err}}
}

func legacyConverter(conf NotifyConfigLegacy) []*NotifyConfig {
	c := &NotifyConfig{
		Host:          conf.Host,
		HostCMD:       conf.HostCMD,
		Port:          conf.Port,
		TLS:           conf.TLS,
		TLSOptions:    conf.TLSOptions,
		Username:      conf.Username,
		UsernameCMD:   conf.UsernameCMD,
		Password:      conf.Password,
		PasswordCMD:   conf.PasswordCMD,
		XOAuth2:       conf.XOAuth2,
		OnNewMail:     conf.OnNewMail,
		OnNewMailPost: conf.OnNewMailPost,
	}
	for _, mailbox := range conf.Boxes {
		c.Boxes = append(c.Boxes, Box{Mailbox: mailbox})
	}
	return []*NotifyConfig{c}
}

func loadConfig(d []byte, debugging bool) ([]*NotifyConfig, error) {
	var config []*NotifyConfig
	err := json.Unmarshal(d, &config)
	if err != nil {
		confErrs := newConfigurationError(err)
		var configLegacy NotifyConfigLegacy
		err = json.Unmarshal(d, &configLegacy)
		if err != nil {
			confErrs.Push(err)
			return nil, confErrs
		} else {
			logrus.Infoln("Legacy configuration format detected")
			config = legacyConverter(configLegacy)
		}
	}

	for i := range config {
		config[i] = retrieveCmd(config[i])
		config[i].Debug = debugging
		if config[i].Alias == "" {
			config[i].Alias = config[i].Username
		}
		for j := range config[i].Boxes {
			config[i].Boxes[j] = setFromConfig(config[i], config[i].Boxes[j])
		}
	}

	return config, nil
}
