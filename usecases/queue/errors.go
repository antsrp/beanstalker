package queue

import "errors"

const (
	msgTubeNotFound = "tube not found"
)

var (
	ErrorTubeNotFound = errors.New(msgTubeNotFound)
)
