package services

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidCode  = status.New(codes.NotFound, "invalid code").Err()
	UserNotFound = status.New(codes.NotFound, "user not found").Err()
	//HelpError     = errors.New("need help")
	//ExistsCode    = errors.New("exists code")
	//AlreadyBinded = errors.New("already binded")
	//MaxAccount    = errors.New("max account")
)
