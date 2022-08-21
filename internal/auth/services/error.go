package services

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	HelpError     = errors.New("need help")
	ExistsCode    = errors.New("exists code")
	AlreadyBinded = errors.New("already binded")
	MaxAccount    = errors.New("max account")
	InvalidCode   = status.New(codes.NotFound, "invalid code").Err()
)
