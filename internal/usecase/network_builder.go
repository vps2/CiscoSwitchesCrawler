package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/vps2/cisco-switches-crawler/internal/domain"
	"github.com/vps2/cisco-switches-crawler/internal/infrastructure/cisco"
	"github.com/vps2/cisco-switches-crawler/pkg/queue"
)

type Client interface {
	Connect(address string, user string, password string) error
	Close() error
	Info() (cisco.ClientInfo, error) //TODO remove direct dependence on the infrastructure layer -> cisco.ClientInfo
}

type IPFilter interface {
	Allow(ip net.IP) bool
}

type Option func(*NetworkBuilder)

func WithShowOutput() Option {
	return func(nb *NetworkBuilder) {
		nb.showOutput = true
	}
}

func WithIPFiltering(filter IPFilter) Option {
	return func(nb *NetworkBuilder) {
		nb.ipFilter = filter
	}
}

type NetworkBuilder struct {
	network     *domain.Network
	ciscoClient Client
	ipFilter    IPFilter
	showOutput  bool
}

func NewNetworkBuilder(cl Client, opts ...Option) *NetworkBuilder {
	nb := &NetworkBuilder{
		network:     domain.NewNetwork(),
		ciscoClient: cl,
	}

	for _, opt := range opts {
		opt(nb)
	}

	return nb
}

func (nb *NetworkBuilder) Build(ctx context.Context, rootSwitchIP string, user string, password string) {
	rootSwitch, _ := domain.NewSwitch(rootSwitchIP)

	queue := queue.New[*domain.Switch]()
	queue.Push(rootSwitch)

	var visited []*domain.Switch

	timerDuration := 3 * time.Second //Switch polling interval. If you do it more often, then the management interface of the switches "falls off".
	timer := time.NewTimer(timerDuration)
	defer timer.Stop()

loop:
	for !queue.IsEmpty() {
		timer.Reset(timerDuration)

		select {
		case <-ctx.Done():
			break loop
		case <-timer.C:
			currSwitch := queue.Pop()
			if inList(visited, currSwitch) {
				continue
			}
			visited = append(visited, currSwitch)

			if err := nb.ciscoClient.Connect(currSwitch.Address(), user, password); err != nil {
				if nb.showOutput {
					log.Println()
				}
				log.Println(err)
				continue
			}
			currSwitchInfo, err := nb.ciscoClient.Info()
			if err != nil {
				if nb.showOutput {
					log.Println()
				}
				log.Println(err)
			}
			nb.ciscoClient.Close()

			if rootSwitch == currSwitch {
				currSwitch.SetName(currSwitchInfo.Name)
				nb.network.AddSwitch(*currSwitch)
			}

			for _, neighborInfo := range currSwitchInfo.Neighbors {
				neighboringSwitch, _ := domain.NewSwitch(neighborInfo.Address)
				neighboringSwitch.SetName(neighborInfo.Name)

				if nb.ipFilter != nil {
					if nb.ipFilter.Allow(net.ParseIP(neighborInfo.Address)) {
						queue.Push(neighboringSwitch)

					} else {
						neighboringSwitch.SetName(neighborInfo.Name + ">>>DISCARDED")
					}
				}

				nb.network.AddSwitch(*neighboringSwitch)
				nb.network.AddLink(*currSwitch, *neighboringSwitch)
			}
		}
	}
}

func (nb *NetworkBuilder) ToJSON() []byte {
	return nb.network.ToJSON()
}

func (nb *NetworkBuilder) ToPrettyJSON() []byte {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, nb.ToJSON(), "", "   ")

	return prettyJSON.Bytes()
}

func inList(container []*domain.Switch, sw *domain.Switch) bool {
	for _, currSwitch := range container {
		if currSwitch.Address() == sw.Address() {
			return true
		}
	}

	return false
}
