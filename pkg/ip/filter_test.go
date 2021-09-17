package ip

import (
	"net"
	"strconv"
	"testing"
)

func TestFilter_Add(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{input: "192.168.1.1", wantErr: false},
		{input: "192.168.1.0/24", wantErr: false},
		{input: "192.168.1.256", wantErr: true},
		{input: "192.168.1.", wantErr: true},
		{input: "", wantErr: true},
		{input: "192.168.1.0/33", wantErr: true},
	}

	filter := NewFilter()
	for i, tt := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			if err := filter.Add(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Filter.Add(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestFilter_Allow(t *testing.T) {
	tests := []struct {
		filter func() *Filter
		input  net.IP
		want   bool
	}{
		{
			filter: func() *Filter {
				return &Filter{allowAnyIfEmpty: true}
			},
			input: net.ParseIP("192.168.1.1"),
			want:  true,
		},
		{
			filter: func() *Filter {
				return &Filter{allowAnyIfEmpty: false}
			},
			input: net.ParseIP("192.168.1.1"),
			want:  false,
		},
		{
			filter: func() *Filter {
				allowedIPs := make(map[string]nothing)
				allowedIPs["192.168.1.1"] = nothing{}

				ip, net, _ := net.ParseCIDR("10.15.1.0/24")
				subnets := []*subnet{{str: ip.String(), ipnet: net}}

				return &Filter{ips: allowedIPs, subnets: subnets, allowAnyIfEmpty: true}
			},
			input: net.ParseIP("192.168.1.1"),
			want:  true,
		},
		{
			filter: func() *Filter {
				allowedIPs := make(map[string]nothing)
				allowedIPs["192.168.1.1"] = nothing{}

				ip, net, _ := net.ParseCIDR("10.15.1.0/24")
				subnets := []*subnet{{str: ip.String(), ipnet: net}}

				return &Filter{ips: allowedIPs, subnets: subnets, allowAnyIfEmpty: true}
			},
			input: net.ParseIP("192.168.1.2"),
			want:  false,
		},
		{
			filter: func() *Filter {
				allowedIPs := make(map[string]nothing)
				allowedIPs["192.168.1.1"] = nothing{}

				ip, net, _ := net.ParseCIDR("10.15.1.0/24")
				subnets := []*subnet{{str: ip.String(), ipnet: net}}

				return &Filter{ips: allowedIPs, subnets: subnets, allowAnyIfEmpty: true}
			},
			input: net.ParseIP("10.15.1.253"),
			want:  true,
		},
		{
			filter: func() *Filter {
				allowedIPs := make(map[string]nothing)
				allowedIPs["192.168.1.1"] = nothing{}

				ip, net, _ := net.ParseCIDR("10.15.1.0/24")
				subnets := []*subnet{{str: ip.String(), ipnet: net}}

				return &Filter{ips: allowedIPs, subnets: subnets, allowAnyIfEmpty: true}
			},
			input: net.ParseIP("10.15.2.1"),
			want:  false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			f := tt.filter()
			if got := f.Allow(tt.input); got != tt.want {
				t.Errorf("Filter.Allow('%s') = %v, want %v", tt.input.String(), got, tt.want)
			}
		})
	}
}
