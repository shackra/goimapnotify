package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// This file is part of goimapnotify
// Copyright (C) 2017-2024	Jorge Javier Araya Navarro

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

type Configuration struct {
	Configurations []NotifyConfig `json:"configurations" yaml:"configurations"`
}

// ConfigurationLegacy holds the old configuration format
type ConfigurationLegacy struct {
	Host              string           `yaml:"host" json:"host"`
	HostCMD           string           `yaml:"hostCMD" json:"hostCMD"`
	Port              int              `yaml:"port" json:"port"`
	TLS               bool             `yaml:"tls" json:"tls"`
	TLSOptions        TLSOptionsStruct `yaml:"tlsOptions" json:"tlsOptions"`
	Username          string           `yaml:"username" json:"username"`
	UsernameCMD       string           `yaml:"usernameCMD" json:"usernameCMD"`
	Password          string           `yaml:"password" json:"password"`
	PasswordCMD       string           `yaml:"passwordCMD" json:"passwordCMD"`
	XOAuth2           bool             `yaml:"xoAuth2" json:"xoAuth2"`
	OnNewMail         string           `yaml:"onNewMail" json:"onNewMail"`
	OnNewMailPost     string           `yaml:"onNewMailPost" json:"onNewMailPost"`
	OnDeletedMail     string           `yaml:"onDeletedMail" json:"onDeletedMail"`
	OnDeletedMailPost string           `yaml:"onDeletedMailPost" json:"onDeletedMailPost"`
	Debug             bool             `yaml:"-" json:"-"`
	Boxes             []string         `yaml:"boxes" json:"boxes"`
}

// NotifyConfig holds the configuration
type NotifyConfig struct {
	Host              string           `yaml:"host" json:"host"`
	HostCMD           string           `yaml:"hostCMD" json:"hostCMD"`
	Port              int              `yaml:"port" json:"port"`
	TLS               bool             `yaml:"tls" json:"tls"`
	TLSOptions        TLSOptionsStruct `yaml:"tlsOptions" json:"tlsOptions"`
	Username          string           `yaml:"username" json:"username"`
	UsernameCMD       string           `yaml:"usernameCMD" json:"usernameCMD"`
	Alias             string           `yaml:"alias" json:"alias"`
	Password          string           `yaml:"password" json:"password"`
	PasswordCMD       string           `yaml:"passwordCMD" json:"passwordCMD"`
	XOAuth2           bool             `yaml:"xoAuth2" json:"xoAuth2"`
	OnNewMail         string           `yaml:"onNewMail" json:"onNewMail"`
	OnNewMailPost     string           `yaml:"onNewMailPost" json:"onNewMailPost"`
	OnDeletedMail     string           `yaml:"onDeletedMail" json:"onDeletedMail"`
	OnDeletedMailPost string           `yaml:"onDeletedMailPost" json:"onDeletedMailPost"`
	Debug             bool             `yaml:"debug" json:"debug"`
	Boxes             []Box            `yaml:"boxes" json:"boxes"`
}

type TLSOptionsStruct struct {
	RejectUnauthorized bool `yaml:"rejectUnauthorized" json:"rejectUnauthorized"`
	STARTTLS           bool `yaml:"starttls" json:"starttls"`
}

/*
Box stores all the necessary info needed to be passed in an
IDLEEvent handler routine, in order to schedule commands and
print informative messages
*/
type Box struct {
	Alias             string    `json:"-" yaml:"-"`
	Mailbox           string    `yaml:"mailbox" json:"mailbox"`
	Reason            EventType `json:"-" yaml:"-"`
	OnNewMail         string    `yaml:"onNewMail" json:"onNewMail"`
	OnNewMailPost     string    `yaml:"onNewMailPost" json:"onNewMailPost"`
	OnDeletedMail     string    `yaml:"onDeletedMail" json:"onDeletedMail"`
	OnDeletedMailPost string    `yaml:"onDeletedMailPost" json:"onDeletedMailPost"`
	ExistingEmail     uint32    `json:"-" yaml:"-"`
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
	c.OnDeletedMail = conf.OnDeletedMail
	c.OnDeletedMailPost = conf.OnDeletedMailPost
	for _, mailbox := range conf.Boxes {
		c.Boxes = append(c.Boxes, Box{Mailbox: mailbox})
	}
	return append(r, c)
}

func loadConfiguration(path string) (Configuration, error) {
	var topConfiguration Configuration
	if err := viper.Unmarshal(&topConfiguration); err != nil {
		return Configuration{}, fmt.Errorf("Can't parse the configuration: %s, error: %v", path, err)
	}

	if topConfiguration.Configurations == nil {
		var legacy ConfigurationLegacy
		if err := viper.UnmarshalExact(&legacy); err != nil {
			return Configuration{}, fmt.Errorf("Can't parse the configuration in 'legacy' format: %s, error: %v", path, err)
		}

		topConfiguration.Configurations = legacyConverter(legacy)
	}

	if len(topConfiguration.Configurations) > 0 && topConfiguration.Configurations[0].Boxes == nil {
		return Configuration{}, fmt.Errorf("configuration file '%s' is empty or have invalid configuration format", path)
	}

	return topConfiguration, nil
}
