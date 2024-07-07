package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// conn, err := net.Dial("tcp", "127.0.0.1:8080")
	// if err != nil {
	// 	// handle error
	// 	log.Fatalln("Error in Dial:", err)
	// }
	// fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	// status, err := bufio.NewReader(conn).ReadString('\n')
	// fmt.Println(status)

	// request := http.NewRequest("GET", )

	post()
	head()
}

func post() {
	fileBytes, err := ioutil.ReadFile("./resource/cscFile.html")
	if err != nil {
		// handle error
		log.Fatalln("Error in ReadFile:", err)
	}
	resp, err := http.Post("http://127.0.0.1:8080/newFile.html", "text/html", bytes.NewReader(fileBytes))
	if err != nil {
		// handle error
		log.Fatalln("Error in Post:", err)
	}

	// read response body
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error in ReadAll:", error)
	}
	// close response body
	resp.Body.Close()

	// print response body
	println(resp.Status)
	println(resp.StatusCode)
	fmt.Println(string(body))
}

func head() {
	resp, err := http.Head("http://127.0.0.1:8080/resource/newCssfile.css")
	if err != nil {
		// handle error
		log.Fatalln("Error in Head:", err)
	}
	// read response body
	body, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error in ReadAll:", error)
	}
	// close response body
	resp.Body.Close()

	// print response body
	println(resp.Status)
	println(resp.StatusCode)
	fmt.Println(string(body))
}
