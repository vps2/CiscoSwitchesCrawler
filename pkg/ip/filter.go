package ip

import (
	"errors"
	"net"
)

type nothing struct{}

type Option func(*Filter)

func AllowAnyIfEmpty(b bool) Option {
	return func(f *Filter) {
		f.allowAnyIfEmpty = b
	}
}

type Filter struct {
	ips             map[string]nothing
	subnets         []*subnet
	allowAnyIfEmpty bool
}

type subnet struct {
	str   string
	ipnet *net.IPNet
}

func NewFilter(opts ...Option) *Filter {
	f := &Filter{
		ips: make(map[string]nothing),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *Filter) Add(addr string) error {
	if ip, net, err := net.ParseCIDR(addr); err == nil { //address with subnet
		if ones, bits := net.Mask.Size(); ones == bits { //containing only one ip? (no bits masked)
			f.addIP(ip.String())
			return nil
		} else {
			f.addSubnet(&subnet{str: ip.String(), ipnet: net})
			return nil
		}
	}
	if ip := net.ParseIP(addr); ip != nil {
		f.addIP(ip.String())
		return nil
	}

	return errors.New("invalid address or subnet")
}

func (f *Filter) Allow(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if f.allowAnyIfEmpty && (len(f.ips) == 0 && len(f.subnets) == 0) {
		return true
	}

	if _, ok := f.ips[ip.String()]; ok {
		return true
	}

	for _, subnet := range f.subnets {
		if subnet.ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

func (f *Filter) addIP(str string) {
	f.ips[str] = nothing{}
}

func (f *Filter) addSubnet(s *subnet) {
	f.subnets = append(f.subnets, s)
}
