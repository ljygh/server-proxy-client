package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Setting struct {
	ServerUrl string `json:"server_url"`
	FilePath  string `json:"file_path"`
	FileType  string `json:"file_type"`
}

func main() {
	filePath := "./resource/cscFile.html"
	serverUrl := "http://localhost:8000/cscFile.html"
	fileType := "text/html"

	if len(os.Args) > 1 {
		settingFilePath := os.Args[1]
		var setting Setting
		getSetting(settingFilePath, &setting)
		filePath = setting.FilePath
		serverUrl = setting.ServerUrl
		fileType = setting.FileType
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		// handle error
		log.Fatalln("Error in ReadFile:", err)
	}
	resp, err := http.Post(serverUrl, fileType, bytes.NewReader(fileBytes))
	if err != nil {
		// handle error
		log.Fatalln("Error in Post:", err)
	}

	// read response body
	body, error := io.ReadAll(resp.Body)
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
