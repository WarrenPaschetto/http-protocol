package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const port = ":42069"
const delimiter = byte('\n')

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.Fatalf("error resolving UDP address: %s\n", err.Error())
	}

	conn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		log.Fatalf("error preparing UDP connection: %s\n", err.Error())
	}
	defer conn.Close()

	b := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		readLine, err := b.ReadString(delimiter)
		if err != nil {
			log.Fatalf("error finding end of line: %s\n", err.Error())
		}

		n, err := conn.Write([]byte(readLine))
		if err != nil {
			log.Fatalf("error writing to UDP connection: %s\n", err.Error())
		}
		log.Printf("%v", n)
	}
}
