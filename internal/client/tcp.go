package client

import (
	"fmt"
	"net"
)

const TCP = "tcp"
const TCPPackageSize = 65535

type TCPClient struct {
	addr string
	conn net.Conn
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

func (c *TCPClient) Receive() {
	buffer := make([]byte, 0, TCPPackageSize)
	for {
		n, err := c.conn.Read(buffer[0:cap(buffer)])
		if err != nil {
			return
		}
		fmt.Println(string(buffer[0:n]))
	}
}

func (c *TCPClient) Send(b []byte) error {
	if _, err := c.conn.Write(b); err != nil {
		return err
	}
	return nil
}

func (c *TCPClient) Stop() {
	c.conn.Close()
}
