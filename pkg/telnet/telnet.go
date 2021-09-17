package telnet

import (
	"fmt"
	"time"

	"github.com/reiver/go-telnet"
)

const defaultWriteTimeout = 200 * time.Millisecond

type Option func(*Client)

// WriteTimeout this is a timeout in msec between write commands
func WriteTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.writeTimeout = timeout
	}
}

// Client wrapper over the telnet package "github.com/reiver/go-telnet"
type Client struct {
	conn         *telnet.Conn
	writeTimeout time.Duration
}

func New(opts ...Option) *Client {
	c := &Client{
		writeTimeout: defaultWriteTimeout,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Connect(address string, port int) error {
	if c.conn != nil {
		return fmt.Errorf("telnet client already connected to: %s", c.conn.RemoteAddr().String())
	}

	conn, err := telnet.DialTo(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return fmt.Errorf("telnet connect: %w", err)
	}

	c.conn = conn

	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err == nil {
			c.conn = nil
		} else {
			return fmt.Errorf("telnet close: %w", err)
		}
	}

	return nil
}

func (c *Client) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func (c *Client) Write(p []byte) (n int, err error) {
	n, err = c.conn.Write(p)
	time.Sleep(c.writeTimeout)
	return
}
