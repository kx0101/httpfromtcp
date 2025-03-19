package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		_, err = conn.Write([]byte(text))
	}
}
