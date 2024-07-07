package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

type Setting struct {
	Port       string `json:"port"`
	ServerUrl  string `json:"server_url"`
	MaxClients int    `json:"max_clients"`
}

// Main function
func main() {
	// Init default port, IP addr, serverUrl and goroutine channels
	const localIP = "127.0.0.1"
	port := "9000"
	serverUrl := "http://localhost:8000"
	maxClients := 2

	// Initialize variables from setting file
	if len(os.Args) > 1 {
		filepath := os.Args[1]
		var setting Setting
		getSetting(filepath, &setting)
		port = setting.Port
		serverUrl = setting.ServerUrl
		maxClients = setting.MaxClients
	}

	ch := make(chan string, maxClients)
	addr := localIP + ":" + port
	println("Server url:", serverUrl)

	// Get TCP addr
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalln("Error in ResolveTCPAddr:", err)
	}

	// Init tcp listener
	tcpLn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln("Error in ListenTCP:", err)
	}
	println("Listen to:", addr)

	// Loop to accept new connections, occupy channel and go to handle connections
	for {
		tcpConn, err := tcpLn.AcceptTCP()
		remoteAddr := tcpConn.RemoteAddr().String()
		println()
		fmt.Println("New client:", remoteAddr)
		if err != nil {
			println("Error in AcceptTCP:", err)
		}
		ch <- remoteAddr
		go handleConnection(*tcpConn, ch, serverUrl)
	}
}

// Get setting from json file
func getSetting(filePath string, setting *Setting) {
	// Open the JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(byteValue, setting)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Handle TCP connection
func handleConnection(conn net.TCPConn, ch chan string, serverUrl string) {
	// Get request
	bufioReader := bufio.NewReader(&conn)
	request, err := http.ReadRequest(bufioReader)
	if err != nil {
		// Read EOF, means tcp disconnects
		if err.Error() == "EOF" {
			println("ReadRequest EOF")
			conn.Close()
			fmt.Printf("Client %v disconnects\n", <-ch)
			return
		}
		fmt.Println("Error in ReadRequest:", err)
		conn.Close()
		fmt.Printf("Client %v disconnects\n", <-ch)
		return
	}

	// Construct response
	response := http.Response{
		Status:        "",
		StatusCode:    200,
		Proto:         "HTTP/1.0",
		Header:        make(http.Header, 0),
		Body:          nil,
		ContentLength: 0,
	}

	// Handle request rather that GET
	if request.Method != "GET" {
		println("Request methods rather than GET not implement")
		sendError(response, "Not implemented", 501, "Request methods rather than GET not implement", conn)
		conn.Close()
		fmt.Printf("Client %v disconnects\n", <-ch)
		return
	}

	// Request to a server and get the response
	url := serverUrl + request.URL.Path
	fmt.Println("GET request for", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in Get:", err)
		sendError(response, "Server internal error", 500, "Error while requesting the server", conn)
		conn.Close()
		fmt.Printf("Client %v disconnects\n", <-ch)
		return
	}

	// Send response to client
	sendResponse(*resp, conn)

	// Disconnect tcp connection
	conn.Close()
	// time.Sleep(10 * time.Second)
	fmt.Printf("Client %v disconnects\n", <-ch)
}

// Send response with a TCP connection
func sendResponse(response http.Response, conn net.TCPConn) {
	buff := bytes.NewBuffer(nil)
	response.Write(buff)
	conn.Write(buff.Bytes())
	println("Send response")
}

// Send error response with TCP connection
func sendError(response http.Response, status string, statusCode int, bodyString string, conn net.TCPConn) {
	response.Status = status
	response.StatusCode = statusCode
	bodyString += "\nError code: " + fmt.Sprint(statusCode)
	response.Body = io.NopCloser(bytes.NewBufferString(bodyString))

	buff := bytes.NewBuffer(nil)
	response.Write(buff)
	conn.Write(buff.Bytes())
	println("Send error:", fmt.Sprint(statusCode))
}
