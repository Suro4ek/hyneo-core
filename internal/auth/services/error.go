package services

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	HelpError   = errors.New("need help")
	InvalidCode = status.Newf(codes.NotFound, "invalid code")
)
