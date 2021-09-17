package cisco

import (
	"bytes"
	"fmt"
	"net"
	"regexp"
	"strings"
)

const (
	newLine = "\n"
	space   = " "

	//different switches have different user name request labels
	txtUsername = "Username:"
	txtLogin    = "Login:"

	txtPassword = "Password:"

	txtAuthenticationFailed = "Authentication failed"
	txtLoginInvalid         = "Login invalid"
	txtBadPasswords         = "Bad passwords"
	txtTimeoutExpired       = "timeout expired"
	txtPrompt               = ">"
	txtMore                 = "--More--"
	txtDeviceSeparator      = "-------------------------"

	cmdShowNeighbors = "sh cdp nei det"
)

const defaultTelnetPort = 23

var re = regexp.MustCompile(`Device ID: (.*?)\r\n.*?\r\n.*?IP address: (.*?)\r\n`)

type Telnet interface {
	Connect(string, int) error
	Close() error
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
}

type Client struct {
	telnet    Telnet
	connected bool
	verbose   bool

	info ClientInfo
}

type Option func(*Client)

func WithVerbose() Option {
	return func(c *Client) {
		c.verbose = true
	}
}

func NewClient(telnet Telnet, opts ...Option) *Client {
	c := &Client{
		telnet: telnet,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Connect(address string, user string, password string) error {
	c.info.Address = ""
	c.info.Name = ""
	c.info.Neighbors = c.info.Neighbors[:0]

	if c.connected {
		c.Close()
	}

	if net.ParseIP(address) == nil {
		return fmt.Errorf("client connect [%v]: wrong ip address", address)
	}
	c.info.Address = address

	if err := c.telnet.Connect(address, defaultTelnetPort); err != nil {
		return fmt.Errorf("client connect [%v]: %w", address, err)
	}

	c.connected = true

	if c.verbose {
		fmt.Println()
	}

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	var serverResponse bytes.Buffer
	for {
		n, err := c.telnet.Read(buffer[:])
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			return fmt.Errorf("client connect [%v]: %w", address, err)
		}

		if c.verbose {
			fmt.Print(string(buffer[:]))
		}

		serverResponse.WriteByte(buffer[0])
		if strings.Contains(serverResponse.String(), txtUsername) ||
			strings.Contains(serverResponse.String(), txtLogin) {
			serverResponse.Reset()
			c.telnet.Write([]byte(user + newLine))
		} else if strings.Contains(serverResponse.String(), txtPassword) {
			serverResponse.Reset()
			c.telnet.Write([]byte(password + newLine))
		} else if strings.Contains(serverResponse.String(), txtAuthenticationFailed) ||
			strings.Contains(serverResponse.String(), txtLoginInvalid) ||
			strings.Contains(serverResponse.String(), txtBadPasswords) {
			serverResponse.Reset()
			c.Close()
			return fmt.Errorf("client connect [%v]: authentication error", address)
		} else if strings.Contains(serverResponse.String(), txtTimeoutExpired) {
			serverResponse.Reset()
			c.Close()
			return fmt.Errorf("client connect [%v]: timeout expired", address)
		} else if strings.HasSuffix(serverResponse.String(), ">") {
			break
		}
	}

	return nil
}

func (c *Client) Close() error {
	if c.connected == false {
		return fmt.Errorf("client close [%v]: connection already closed", c.info.Address)
	}

	c.connected = false

	return c.telnet.Close()
}

func (c *Client) Info() (ClientInfo, error) {
	if !c.connected {
		return c.info, fmt.Errorf("client info [%v]: connection closed", c.info.Address)
	}

	if _, err := c.telnet.Write([]byte(newLine)); err != nil {
		return c.info, fmt.Errorf("client info [%v]: %w", c.info.Address, err)
	}

	var serverResponse bytes.Buffer
	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	var isCmdStarted = false
	for {
		n, err := c.telnet.Read(buffer[:])
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			return c.info, fmt.Errorf("client info [%v]: %w", c.info.Address, err)
		}

		if c.verbose {
			fmt.Print(string(buffer[:]))
		}

		serverResponse.WriteByte(buffer[0])
		if strings.HasSuffix(serverResponse.String(), txtPrompt) {
			if isCmdStarted == false {
				c.info.Name = strings.TrimSpace(strings.TrimRight(serverResponse.String(), txtPrompt))
				serverResponse.Reset()

				c.telnet.Write([]byte(cmdShowNeighbors + newLine))
				isCmdStarted = true
			} else {
				input := strings.Replace(serverResponse.String()+newLine, cmdShowNeighbors, "", -1)
				input = strings.TrimSpace(input)
				input = strings.Replace(input, txtDeviceSeparator+newLine, "", 1)
				serverResponse.Reset()

				c.info.Neighbors = parseInput(input)

				break
			}
		} else if strings.Contains(serverResponse.String(), txtMore) {
			serverResponse.Truncate(serverResponse.Len() - len(txtMore))
			c.telnet.Write([]byte(space))
		}
	}

	return c.info, nil
}

func parseInput(in string) []ClientInfo {
	var neighbors []ClientInfo
	tokens := strings.Split(in, txtDeviceSeparator)
	for _, t := range tokens {
		res := re.FindStringSubmatch(t)
		if res != nil {
			neighbors = append(neighbors, ClientInfo{Name: res[1], Address: res[2]})
		}
	}

	return neighbors
}
