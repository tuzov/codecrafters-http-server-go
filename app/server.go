package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
	var responseBody []byte

	switch request.Method {
	case "GET":
		response, responseBody = HandleGet(request)
	case "POST":
		response, responseBody = HandlePost(request)
	default:
		response = "HTTP/1.1 405 Method Not Allowed\r\n\r\n"
		responseBody = nil
	}
	_, err = conn.Write([]byte(response))
	_, err = conn.Write(responseBody)
	if err != nil {
		fmt.Println("Error writing to connection: ", err)
		return
	}
}

func HandlePost(request *http.Request) (string, []byte) {
	filename := "/tmp/data/codecrafters.io/http-server-tester/" + strings.Split(request.URL.Path, "/")[2]
	body, _ := io.ReadAll(request.Body)
	err := os.WriteFile(filename, body, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("HTTP/1.1 201 Created\r\n\r\n"), nil
}

func HandleGet(request *http.Request) (string, []byte) {
	switch path := request.URL.Path; {
	case strings.HasPrefix(path, "/echo/"):
		str := strings.Split(path, "/")[2]
		compression := request.Header.Get("Accept-Encoding")
		if strings.Contains(compression, "gzip") {
			//compressed := gzipCompress([]byte(str))
			var buffer bytes.Buffer
			w := gzip.NewWriter(&buffer)
			w.Write([]byte(str))
			w.Close()
			compressed := buffer.Bytes()
			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: %d\r\n\r\n", len(compressed))
			return response, compressed
		} else {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len([]byte(str))), []byte(str)
		}

	case path == "/user-agent":
		return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n", len(request.UserAgent())), []byte(request.UserAgent())
	case strings.HasPrefix(path, "/files/"):
		filename := "/tmp/data/codecrafters.io/http-server-tester/" + strings.Split(path, "/")[2]
		dat, err := os.ReadFile(filename)
		fmt.Println(dat, string(dat))
		if err != nil {
			return "HTTP/1.1 404 Not Found\r\n\r\n", nil
		} else {
			return fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", len(dat)), dat
		}
	case path == "/":
		return "HTTP/1.1 200 OK\r\n\r\n", nil
	default:
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}
}
