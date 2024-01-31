package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/oleksandrcherevkov/go-tcp-messenger/internal/server"
)

const (
	port    = ":40000"
	maxConn = 20
)

func main() {
	server := server.NewTCP(port, maxConn)
	go server.Run()
	go server.Broadcast()
	fmt.Println("Click to Stop")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
