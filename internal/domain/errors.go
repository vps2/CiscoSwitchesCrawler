package domain

import "errors"

var (
	ErrInvalidSwitchIPAddress = errors.New("wrong switch ip address")

	ErrEmptySwitchAddress = errors.New("empty switch IP address")
	ErrSwitchNotInNetwork = errors.New("the switch has not been added to the network")
	ErrLink               = errors.New("attempt to create a link with yourself")
)
