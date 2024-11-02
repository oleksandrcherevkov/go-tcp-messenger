package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/oleksandrcherevkov/go-tcp-messenger/internal/client"
)

const (
	serverPort = ":40000"
)

func main() {
	client, err := client.NewTCP(serverPort)
	if err != nil {
		log.Fatalf("Failed connecting to %v. %v", serverPort, err)
	}
	requests := client.Receive()
	go func() {
		for request := range requests {
			fmt.Println(string(request))
		}
	}()
	defer client.Stop()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "q" {
			break
		}
		err = client.Send(scanner.Bytes())
		if err != nil {
			log.Fatalf("Failed sending message %v. %v", scanner.Text(), err)
		}
	}

	err = scanner.Err()
	if err != nil {
		log.Fatalf("Failed reading from console. %v", err)
	}
}
