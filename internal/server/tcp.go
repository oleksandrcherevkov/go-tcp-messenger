package server

import (
	"log"
	"log/slog"
	"net"
)

const TCP = "tcp"
const TCPPackageSize = 65535

type TCPServer struct {
	port     string
	maxConn  uint8
	users    *UsersRepository
	listener net.Listener
	messages chan *Message
	stop     chan interface{}
}

func NewTCP(port string, maxConnections uint8) *TCPServer {
	return &TCPServer{
		port:     port,
		maxConn:  maxConnections,
		messages: make(chan *Message),
		users:    NewUsersRepo(),
	}
}

func (s *TCPServer) Run() {
	listener, err := net.Listen(TCP, s.port)
	if err != nil {
		log.Fatalf("Error while creating Listener. %v", err)
	}
	if listener == nil {
		log.Fatalf("Created nil listener.")
	}
	defer listener.Close()
	s.listener = listener

	for {
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(c net.Conn) {
	user := NewUser(c)
	defer user.Leave()
	slog.Info("Connected: \n", "Addr", user.Addr)
	s.users.AddUser(user)
	defer s.users.RemoveUser(user.Addr)
	buffer := make([]byte, 0, TCPPackageSize)

	for {
		n, err := user.conn.Read(buffer[0:cap(buffer)])
		if err != nil {
			break
		}
		message := NewMessage(buffer[0:n])
		slog.Info("Message received: \n", "Addr", user.Addr, "Message", string(message.text))
		s.messages <- message
	}
	slog.Info("Disconnected: \n", "Addr", user.Addr)
}

func (s *TCPServer) Broadcast() {
	for m := range s.messages {
		s.sendToAll(m)
	}
}

func (s *TCPServer) sendToAll(m *Message) {
	s.users.IterateUsers(func(u *User) {
		u.Send(m.text)
	})
}
