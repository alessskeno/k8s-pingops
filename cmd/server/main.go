// cmd/server/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"k8s-pingops/pkg/telnet"
)

type TelnetRequest struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type TelnetResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		var req TelnetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
