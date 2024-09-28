// cmd/server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"k8s-pingops/pkg/dns"
	"k8s-pingops/pkg/httpclient"
	"k8s-pingops/pkg/telnet"
)

type Request struct {
	Host     string            `json:"host"`
	Port     int               `json:"port,omitempty"`
	Function string            `json:"function"`
	URL      string            `json:"url,omitempty"`
	Method   string            `json:"method,omitempty"`
	Body     string            `json:"body,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

type TelnetResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type DNSResponse struct {
	Success bool     `json:"success"`
	IPs     []string `json:"ips,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type HTTPResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
	Body    string `json:"body,omitempty"`
	Error   string `json:"error,omitempty"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch req.Function {
		case "telnet":
			err := telnet.CheckTelnetConnection(req.Host, req.Port)
			var res TelnetResponse
			if err != nil {
				res = TelnetResponse{Success: false, Error: err.Error()}
			} else {
				res = TelnetResponse{Success: true}
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "dns":
			ips, err := dns.ResolveDNS(req.Host)
			var res DNSResponse
			if err != nil {
				res = DNSResponse{Success: false, Error: err.Error()}
			} else {
				res = DNSResponse{Success: true, IPs: ips}
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case "http":
			res := httpclient.MakeHTTPRequest(req.URL, req.Method, req.Body, req.Headers)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		default:
			http.Error(w, "Invalid function", http.StatusBadRequest)
		}
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
