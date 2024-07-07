# Server-proxy-client

## Environment
Install golang as shown in: https://go.dev/doc/install

## Running

### 1. Running the server

To run the server enter the following command in the command window:
```go run server.go [port_number] [max_client_number]```  
Port nubmer is optional and default port number is 8000.
Max client number is optional and the default one is 2.


### 2. Running the proxy

To run the proxy you need two arguments one for the portnumber and one for the URL:
-go run proxy.go 8080 https://google.com/

However the arguments are not neccesary to run the proxy

### 3. Running the client

## Implementation

### 1. Server
It builds a tcp listener with ip address and port and loops to receive new connections from clients. It uses channels to limit the max number of clients the server serves at the same time. It receives http request from clients and sends http response back.

