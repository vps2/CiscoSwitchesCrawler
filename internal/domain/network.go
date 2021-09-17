package domain

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/vps2/cisco-switches-crawler/pkg/set"
)

var switchRegexp = regexp.MustCompile(`\{(.*? .*?)\}`)

// Network - network of switches
type Network struct {
	graph map[Switch]*set.Set[Switch]
}

func NewNetwork() *Network {
	return &Network{
		graph: make(map[Switch]*set.Set[Switch]),
	}
}

func (n *Network) AddSwitch(s Switch) error {
	if s.Address() == "" {
		return fmt.Errorf("network add switch [%s]: %w", s.Address(), ErrEmptySwitchAddress)
	}

	if _, ok := n.graph[s]; !ok {
		n.graph[s] = set.New[Switch]()
	}

	return nil
}

func (n *Network) AddLink(fromSwitch, toSwitch Switch) error {
	fromSwitchNeighbors, fromSwitchFound := n.graph[fromSwitch]
	toSwitchNeighbors, toSwitchFound := n.graph[toSwitch]

	if !fromSwitchFound {
		return fmt.Errorf("network add link to [%s]: %w", fromSwitch.Address(), ErrSwitchNotInNetwork)
	}
	if !toSwitchFound {
		return fmt.Errorf("network add link to [%s]: %w", toSwitch.Address(), ErrSwitchNotInNetwork)
	}

	if fromSwitch == toSwitch {
		return fmt.Errorf("network add link to [%s]: %w", fromSwitch.Address(), ErrLink)
	}

	fromSwitchNeighbors.Add(toSwitch)
	toSwitchNeighbors.Add(fromSwitch)

	return nil
}

func (n *Network) NeighborsOf(sw Switch) ([]Switch, error) {
	neighbors, ok := n.graph[sw]
	if !ok {
		return []Switch{}, fmt.Errorf("network show neighbors [%s]: %w", sw.Address(), ErrSwitchNotInNetwork)
	}

	return neighbors.ToSlice(), nil
}

func (n *Network) Len() int {
	return len(n.graph)
}

func (n *Network) ToJSON() []byte {
	graphSize := len(n.graph)
	i := 0

	var buf bytes.Buffer
	buf.WriteString(`{"network": [`)
	for k, v := range n.graph {
		buf.WriteString("{")
		buf.WriteString(fmt.Sprintf(`"name": "%s", "address": "%s"`, k.Name(), k.Address()))

		submatches := switchRegexp.FindAllStringSubmatch(v.String(), -1)
		size := len(submatches)
		if size > 0 {
			buf.WriteString(`, "neighbors": [`)
			for j, submatch := range submatches {
				tokens := strings.Split(submatch[1], " ")
				buf.WriteString(fmt.Sprintf(`{"name": "%s", "address": "%s"}`, tokens[0], tokens[1]))
				if j < size-1 {
					buf.WriteString(",")
				}
			}
			buf.WriteString("]")
		}
		buf.WriteString("}")

		if i < graphSize-1 {
			buf.WriteString(",")
		}
		i++
	}
	buf.WriteString("]}")

	return buf.Bytes()
}

func (n *Network) String() string {
	return string(n.ToJSON())
}
