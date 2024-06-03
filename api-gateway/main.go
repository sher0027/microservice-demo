package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func proxyRequest(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/")
	parts := strings.SplitN(trimmedPath, "/", 2)
	if len(parts) == 0 {
		http.Error(w, "invalid service name", http.StatusBadRequest)
		return
	}
	serviceName := parts[0]
	log.Printf("Received request for service: %s", serviceName)

	resp, err := http.Get("http://localhost:8761/service?name=" + serviceName)
	if err != nil {
		log.Printf("Error getting service address: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Service %s not available, status code: %d", serviceName, resp.StatusCode)
		http.Error(w, "service not available", http.StatusInternalServerError)
		return
	}

	addr, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading service address response body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Service address for %s: %s", serviceName, string(addr))

	proxyURL := "http://" + string(addr) + r.URL.Path
	if r.URL.RawQuery != "" {
		proxyURL += "?" + r.URL.RawQuery
	}
	log.Printf("Forwarding request to: %s", proxyURL)

	newReq, err := http.NewRequest(r.Method, proxyURL, r.Body)
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newReq.Header = r.Header

	client := &http.Client{}
	serviceResp, err := client.Do(newReq)
	if err != nil {
		log.Printf("Error forwarding request to service %s: %v", serviceName, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer serviceResp.Body.Close()

	if serviceResp.StatusCode != http.StatusOK && serviceResp.StatusCode != http.StatusCreated {
		log.Printf("Service %s responded with status code: %d", serviceName, serviceResp.StatusCode)
		http.Error(w, "service response error", serviceResp.StatusCode)
		return
	}

	log.Printf("Successfully forwarded request to service %s", serviceName)

	io.Copy(w, serviceResp.Body)
}

func main() {
	http.HandleFunc("/", proxyRequest)
	log.Println("API Gateway is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
