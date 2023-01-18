package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserNotFound      = status.New(codes.NotFound, "user not found").Err()
	Fault             = status.New(codes.Unknown, "fault").Err()
	IncorrectPassword = status.New(codes.Unauthenticated, "incorrect password").Err()
	AccountsLimit     = status.New(codes.ResourceExhausted, "accounts limit").Err()
)
