package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

// Main function
func main() {
	var port = "8000"
	var maxClients = 2

	// Get port and max number of clients from arguments
	if len(os.Args) > 1 {
		_, err := strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln("Error: port_number should be an integer.")
			return
		} else {
			port = os.Args[1]
		}
	}
	if len(os.Args) > 2 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			// handle the error here
			log.Fatalln("Error: max_client_number should be an integer.")
			return
		} else {
			maxClients = i
		}
	}

	// Make channels for concurrency
	ch := make(chan string, maxClients)

	// Init TCP address
	addr := "127.0.0.1:" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatalln("Error in ResolveTCPAddr:", err)
	}

	// Init TCP listener
	tcpLn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalln("Error in ListenTCP:", err)
	}
	fmt.Printf("Listen to %v:%v\n", tcpAddr.IP, tcpAddr.Port)
	println("Max number of clients:", maxClients)

	// Loop to accept new connection
	for {
		tcpConn, err := tcpLn.AcceptTCP()
		if err != nil {
			println("Error in AcceptTCP:", err)
		}
		remoteAddr := tcpConn.RemoteAddr().String()
		println()
		println("New client:", remoteAddr)
		ch <- remoteAddr
		go handleConnection(*tcpConn, ch)
	}
}

// Handle TCP connection, get http request and send response back
func handleConnection(conn net.TCPConn, ch chan string) {
	// Read request
	bufioReader := bufio.NewReader(&conn)
	request, err := http.ReadRequest(bufioReader)
	if err != nil {
		conn.Close()
		println("Error while reading request")
		log.Fatalf("Client %v disconnets\n", <-ch)
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

	// Handle request according to its method
	if request.Method == "GET" {
		getHandler(response, *request, conn)
	} else if request.Method == "POST" {
		postHandler(response, *request, conn)
	} else {
		println("Request Method is not supported")
		sendError(response, "Not Implemented", 501, "Request Method is not supported\nError code: 501", conn)
	}

	// disconnect tcp connection
	conn.Close()
	fmt.Printf("Client %v disconnets\n", <-ch)
}

// Handle GET request
func getHandler(response http.Response, request http.Request, conn net.TCPConn) {
	// get requested file path
	filePath := request.URL.Path
	filePath = "." + filePath
	fmt.Println("GET request for:", filePath)

	// check file type
	fileType := getFileType(filePath)
	if !validFileType(fileType) {
		println("File type not supported")
		sendError(response, "Bad request", 400, "File type not supported", conn)
		return
	}

	// check if file exists
	if _, err := os.Stat(filePath); err != nil {
		println("Requested file not exist")
		sendError(response, "Not found", 404, "Requested file not exist", conn)
		return
	}

	// edit response status
	response.Status = "Status OK"
	response.StatusCode = 200

	// edit response header
	switch fileType {
	case "jpg":
		response.Header.Add("Content-Type", "image/jpg")
	case "jpeg":
		response.Header.Add("Content-Type", "image/jpeg")
	case "html":
		response.Header.Add("Content-Type", "text/html")
	case "txt":
		response.Header.Add("Content-Type", "text/plain")
	case "css":
		response.Header.Add("Content-Type", "text/css")
	case "gif":
		response.Header.Add("Content-Type", "image/gif")
	default:
	}

	// read file and write to response body
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		println("Server error in ReadFile:", err)
		sendError(response, "Internal server error", 500, "Server error in ReadFile:", conn)
	}
	response.Body = io.NopCloser(bytes.NewBuffer(fileBytes))
	response.ContentLength = int64(len(fileBytes))
	// send response
	sendResponse(response, conn)
	// time.Sleep(10 * time.Second) // Used to test goroutines
}

// Handle POST request
func postHandler(response http.Response, request http.Request, conn net.TCPConn) {
	// get post file path
	filePath := request.URL.Path
	filePath = "." + filePath
	fmt.Println("POST request to:", filePath)

	// check file name
	filename := path.Base(filePath)
	if filename == "." || filename == "/" {
		println("No filename in the post request")
		sendError(response, "Bad request", 400, "No filename in the post request", conn)
		return
	}

	// check file type
	fileType := getFileType(filePath)
	if !validFileType(fileType) {
		println("File type not supported")
		sendError(response, "Bad request", 400, "File type not supported", conn)
		return
	}

	// check if directory of file exists
	fileDir := path.Dir(filePath)
	if _, err := os.Stat(fileDir); err != nil {
		println("Post folder not exist")
		sendError(response, "Bad request", 400, "Post folder not exist", conn)
		return
	}

	// Get content from request.Body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		println("Server error while reading request's body", err)
		sendError(response, "Server internal error", 500, "Server error while reading request's body", conn)
		return
	}

	// Create a file to save the content according to request's url
	println(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		println("Server error while creating a new file", err)
		sendError(response, "Server internal error", 500, "Server error while creating a new file", conn)
		return
	}

	// write into file
	_, err = file.Write(body)
	if err != nil {
		println("Server error while writing request's body to the new file", err)
		sendError(response, "Server internal error", 500, "Server error while writing request's body to the new file", conn)
		return
	}

	// close file
	err = file.Close()
	if err != nil {
		println("Server error while closing a new file", err)
		sendError(response, "Server internal error", 500, "Server error while closing a new file", conn)
		return
	}
	println("Save file to:", filePath)

	// Construct response to show information and send it
	response.Status = "Status OK"
	response.StatusCode = 200
	response.Body = io.NopCloser(bytes.NewBufferString("Upload successfully\n"))
	response.ContentLength = int64(len("Upload successfully\n"))
	sendResponse(response, conn)
	// time.Sleep(10 * time.Second) // Used to test goroutines
}

// Send response with TCP connection
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

// Get file type according to extension
func getFileType(filePath string) string {
	var extension = filepath.Ext(filePath)
	var fileType string
	if len(extension) > 0 {
		fileType = extension[1:]
	} else {
		fileType = extension
	}
	return fileType
}

// Check if file type in the list of file types that the server can handle
func validFileType(fType string) bool {
	switch fType {
	case "html", "txt", "gif", "jpeg", "jpg", "css":
		return true
	default:
		return false
	}
}
