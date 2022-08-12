package services

import "errors"

var (
	HelpError   = errors.New("need help")
	InvalidCode = errors.New("invalid code")
)
