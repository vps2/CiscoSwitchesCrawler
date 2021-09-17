package domain_test

import (
	"reflect"
	"testing"

	"github.com/vps2/cisco-switches-crawler/internal/domain"
)

func TestNeighborsOf(t *testing.T) {
	sw1, _ := domain.NewSwitch("192.168.1.1")
	sw1.SetName("sw1")
	sw2, _ := domain.NewSwitch("192.168.1.2")
	sw2.SetName("sw2")
	sw3, _ := domain.NewSwitch("192.168.1.3")
	sw3.SetName("sw3")
	sw4, _ := domain.NewSwitch("192.168.1.4")
	sw4.SetName("sw4")
	sw5, _ := domain.NewSwitch("192.168.1.5")
	sw5.SetName("sw5")

	network := domain.NewNetwork()
	network.AddSwitch(*sw1)
	network.AddSwitch(*sw2)
	network.AddSwitch(*sw3)
	network.AddSwitch(*sw4)
	network.AddSwitch(*sw5)

	network.AddLink(*sw1, *sw2)
	network.AddLink(*sw2, *sw3)
	network.AddLink(*sw2, *sw4)
	network.AddLink(*sw3, *sw5)

	tests := []struct {
		name   string
		in     domain.Switch
		expect []domain.Switch
	}{
		{
			"sw1 neighbors",
			*sw1,
			[]domain.Switch{*sw2},
		},
		{
			"sw2 neighbors",
			*sw2,
			[]domain.Switch{*sw1, *sw3, *sw4},
		},
		{
			"sw3 neighbors",
			*sw3,
			[]domain.Switch{*sw2, *sw5},
		},
		{
			"sw4 neighbors",
			*sw4,
			[]domain.Switch{*sw2},
		},
		{
			"sw5 neighbors",
			*sw5,
			[]domain.Switch{*sw3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := network.NeighborsOf(tt.in); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("Network.NeighborsOf() = %v, want %v", got, tt.expect)
			}
		})
	}
}
