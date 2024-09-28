package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPResponse represents the response from an HTTP request
type HTTPResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
	Body    string `json:"body,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MakeHTTPRequest handles sending HTTP requests with specified method, URL, body, and headers.
func MakeHTTPRequest(url, method, body string, headers map[string]string) HTTPResponse {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return HTTPResponse{Success: false, Error: err.Error()}
	}

	// Add headers if any are provided
	for key, value := range headers {
		req.Header[key] = []string{value}
	}

	if _, exists := req.Header["Content-Type"]; !exists {
		req.Header.Set("Content-Type", "application/json")
	}

	// Log headers for debugging
	fmt.Println("Headers:")
	for key, value := range req.Header {
		fmt.Printf("%s: %s\n", key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return HTTPResponse{Success: false, Error: err.Error()}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return HTTPResponse{Success: false, Error: err.Error()}
	}

	return HTTPResponse{
		Success: true,
		Status:  resp.Status,
		Body:    string(respBody),
	}
}
