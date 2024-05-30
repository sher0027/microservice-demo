package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type ServiceRegistry struct {
	services map[string]string
	mu       sync.RWMutex
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

func (sr *ServiceRegistry) RegisterService(w http.ResponseWriter, r *http.Request) {
	var service struct {
		Name string `json:"name"`
		Addr string `json:"addr"`
	}
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sr.mu.Lock()
	sr.services[service.Name] = service.Addr
	sr.mu.Unlock()
	log.Printf("Registered service: %s at %s", service.Name, service.Addr)
	w.WriteHeader(http.StatusOK)
}

func (sr *ServiceRegistry) GetService(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	sr.mu.RLock()
	addr, ok := sr.services[name]
	sr.mu.RUnlock()
	if !ok {
		log.Printf("Service not found: %s", name)
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	log.Printf("Service found: %s at %s", name, addr)
	w.Write([]byte(addr))
}

func main() {
	sr := NewServiceRegistry()

	http.HandleFunc("/register", sr.RegisterService)
	http.HandleFunc("/service", sr.GetService)

	log.Println("Service Registry is running on port 8761")
	log.Fatal(http.ListenAndServe(":8761", nil))
}
