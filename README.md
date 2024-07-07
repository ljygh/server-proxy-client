# Server-proxy-client
This project implements server, proxy and client with golang. Client can request to server directly or with the help of proxy. Communications among them are based on tcp and http. They can be implemented with net/http library only. Tcp is used in order to explore tcp and goroutine. 

## Environment
Install golang as shown in: https://go.dev/doc/install.

## Running

### 1. Running the server

To run the server, enter the following command in the command window:
```go run server.go [path to server setting file]```  
Default port number is 8000.  
Default max client number is 2.  
./server/server.json is an example of server setting file.

### 2. Running the proxy

To run the proxy enter the following command in the command window:
```go run proxy.go [path to proxy setting file]```
Default port number: "9000"  
Default server url: "http://localhost:8000"  
Default max client number: 2.   

### 3. Running the client
1. Test server and proxy with web browser. It can GET files and show them with server or proxy.
2. Test server with POST in this way:
```go run client.go [path to client setting file]```
Default file path: "./resource/cscFile.html"
Default server Url := "http://localhost:8000/cscFile.html"
Default file type := "text/html"

## Implementation

### 1. Server
It builds a tcp listener with ip address and port and loops to receive new connections from clients. It uses channels to limit the max number of clients the server serves at the same time. It receives http request from clients and sends http response back. The server can handle both GET and POST. For GET, it send the requested file in the response. For POST, it saves request body to a file specified by the request.

### 2. Proxy
It only supports GET request. It will send the same request to the server and send the response of server to client.

### 3. Client
Implement post with http.

