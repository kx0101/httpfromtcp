package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("message.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer file.Close()

	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		currentLine := strings.Builder{}
		buf := make([]byte, 8)

		for {
			n, err := f.Read(buf)
			if err != nil {
				fmt.Errorf("Error: %v", err)
				return
			}

			if err == io.EOF {
				return
			}

			parts := strings.Split(string(buf[:n]), "\n")
			for i, part := range parts {
				if i == len(parts)-1 {
					currentLine.WriteString(part)
					continue
				}

				out <- currentLine.String() + part
				currentLine.Reset()
			}
		}
	}()

	return out
}
