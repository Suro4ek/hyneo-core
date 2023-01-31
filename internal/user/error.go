package user

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	NotFound = status.New(codes.NotFound, "user not found").Err()
	Fault    = status.New(codes.Unknown, "fault").Err()
)
