package domain

import (
	"fmt"
	"net"
)

// Switch provides information about the switch
type Switch struct {
	name    string
	address string
}

func NewSwitch(address string) (*Switch, error) {
	if net.ParseIP(address) == nil {
		return nil, fmt.Errorf("switch new [%s]: %w", address, ErrInvalidSwitchIPAddress)
	}

	return &Switch{
		address: address,
	}, nil
}

func (s *Switch) SetName(name string) {
	s.name = name
}

func (s *Switch) Name() string {
	return s.name
}

func (s *Switch) Address() string {
	return s.address
}

func (s *Switch) String() string {
	return fmt.Sprintf("Switch {Name: %s, Address: %s}", s.name, s.address)
}
