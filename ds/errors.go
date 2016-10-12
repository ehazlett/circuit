package ds

import "errors"

var (
	ErrNetworkDoesNotExist = errors.New("network does not exist")
	ErrServiceDoesNotExist = errors.New("service does not exist")
)
