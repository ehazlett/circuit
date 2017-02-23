package ds

import "errors"

var (
	ErrNetworkDoesNotExist     = errors.New("network does not exist")
	ErrServiceDoesNotExist     = errors.New("service does not exist")
	ErrNetworkPeerDoesNotExist = errors.New("peer does not exist for that network")
)
