package client

import (
	"net"
)

const (
	TCP            = "tcp"
	TCPPackageSize = 65535
)

type TCPClient struct {
	addr     string
	conn     net.Conn
	requests chan []byte
}

func NewTCP(addr string) (*TCPClient, error) {
	conn, err := net.Dial(TCP, addr)
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		addr: addr,
		conn: conn,
	}, nil
}

func (c *TCPClient) Receive() chan []byte {
	buffer := make([]byte, 0, TCPPackageSize)
	requests := make(chan []byte)
	c.requests = requests
	go func() {
		for {
			n, err := c.conn.Read(buffer[0:cap(buffer)])
			if err != nil {
				return
			}
			requests <- buffer[0:n]
		}
	}()
	return requests
}

func (c *TCPClient) Send(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *TCPClient) Stop() {
	close(c.requests)
	c.conn.Close()
}
