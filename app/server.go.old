package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Server is listening on port 4221")

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	HandleClient(conn)

}

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// Read data
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection: ", err)
		return
	}
	request := string(buf[:n])
	//fmt.Println("Received data:\n", request)
	path := strings.Fields(request)[1]

	if path == "/" {
		_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing 200 to connection: ", err)
			return
		}
	} else if strings.Contains(path, "/echo/") {
		str := strings.Split(path, "/")[2]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len([]byte(str)), str)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing 200 to connection: ", err)
			return
		}
	} else if strings.Contains(path, "/user-agent") {
		parts := strings.Split(request, "\r\n")
		for _, j := range parts {
			if strings.Contains(j, "User-Agent:") {
				headerPayload := strings.Split(j, ": ")
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len([]byte(headerPayload[1])), headerPayload[1])
				_, err = conn.Write([]byte(response))
				if err != nil {
					fmt.Println("Error writing 200 to connection: ", err)
					return
				}
			}
		}

	} else {
		_, err = conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		if err != nil {
			fmt.Println("Error writing 404 to connection: ", err)
			return
		}
	}
}
