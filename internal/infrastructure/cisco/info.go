package cisco

import (
	"fmt"
	"strings"
)

type ClientInfo struct {
	Name      string
	Address   string
	Neighbors []ClientInfo
}

func (ci ClientInfo) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ClientInfo {Name: %s, Address: %s", ci.Name, ci.Address))
	if len(ci.Neighbors) == 0 {
		sb.WriteString("}")
	} else {
		sb.WriteString(fmt.Sprintf(", Neighbors: %s}", ci.Neighbors))
	}

	return sb.String()
}
