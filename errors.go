package main

import (
	"errors"
)

var (
	CannotCheckSupportedAuthErr = errors.New("there was an error while checking supported authentication mechanism")
	TokenAuthNotSupportedErr    = errors.New("XOAUTH2 and OAUTHBEARER are not supported by the server")
)
