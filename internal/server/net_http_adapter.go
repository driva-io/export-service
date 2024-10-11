package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type NetHttpClient struct{}

var _ HttpClient = (*NetHttpClient)(nil)

func (s *NetHttpClient) Get(url string, headers map[string]interface{}) (response HttpResponse, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		return
	}

	err = json.Unmarshal(data, &response.Body)
	if err != nil {
		log.Printf("Failed to unmarshal body response: %v\n", err)
		return
	}

	response.StatusCode = resp.StatusCode
	response.Headers = make(map[string]string)
	for key, value := range resp.Header {
		response.Headers[key] = value[0]
	}
	return
}

func (s *NetHttpClient) Post(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Failed to marshal body: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value.(string))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		return
	}

	err = json.Unmarshal(data, &response.Body)
	if err != nil {
		log.Printf("Failed to unmarshal body response: %v\n", err)
		return
	}

	response.StatusCode = resp.StatusCode
	response.Headers = make(map[string]string)
	for key, value := range resp.Header {
		response.Headers[key] = value[0]
	}

	return
}

func (s *NetHttpClient) Patch(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	return
}
func (s *NetHttpClient) Update(url string, body any, headers map[string]interface{}) (response HttpResponse, err error) {
	return
}
func (s *NetHttpClient) Delete(url string, headers map[string]interface{}) (response HttpResponse, err error) {
	return
}
