package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go HandleClient(conn)
	}
}

func HandleClient(conn net.Conn) {
	defer conn.Close()
	// Read data
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request. ", err.Error())
		return
	}

	var response string

	switch request.Method {
	case "GET":
		response = HandleGet(request)
	case "POST":
		response = HandlePost(request)
	default:
		response = "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
	}
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err)
		return
	}
}

func HandlePost(request *http.Request) string {
	filename := "/tmp/data/codecrafters.io/http-server-tester/" + strings.Split(request.URL.Path, "/")[2]
	body, _ := io.ReadAll(request.Body)
	err := os.WriteFile(filename, body, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("HTTP/1.1 201 Created\r\n\r\n")
}

func HandleGet(request *http.Request) string {
	switch path := request.URL.Path; {
	case strings.HasPrefix(path, "/echo/"):
		str := strings.Split(path, "/")[2]
		compression := request.Header.Get("Accept-Encoding")
		if compression == "gzip" {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n%s", len([]byte(str)), str)
		} else {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len([]byte(str)), str)
		}

	case path == "/user-agent":
		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(request.UserAgent()), request.UserAgent())
	case strings.HasPrefix(path, "/files/"):
		filename := "/tmp/data/codecrafters.io/http-server-tester/" + strings.Split(path, "/")[2]
		dat, err := os.ReadFile(filename)
		if err != nil {
			return "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(dat), dat)
		}
	case path == "/":
		return "HTTP/1.1 200 OK\r\n\r\n"
	default:
		return "HTTP/1.1 404 Not Found\r\n\r\n"
	}
}
