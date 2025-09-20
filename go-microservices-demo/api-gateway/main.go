package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

const gatewayPort = 8080

func main() {
	http.HandleFunc("/", routeRequest)

	log.Printf("API Gateway starting on port %d...", gatewayPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", gatewayPort), nil); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

func routeRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Gateway received request for: %s", r.URL.Path)

	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathSegments) < 3 || pathSegments[0] != "api" {
		http.Error(w, "Invalid request path", http.StatusBadRequest)
		return
	}
	serviceName := pathSegments[1] + "-service"

	serviceURL, err := discoverService(serviceName)
	if err != nil {
		log.Printf("Error discovering service '%s': %v", serviceName, err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	log.Printf("Discovered '%s' at %s", serviceName, serviceURL)

	proxy := httputil.NewSingleHostReverseProxy(serviceURL)

	r.URL.Path = "/" + strings.Join(pathSegments[1:], "/")
	log.Printf("Forwarding request to: %s%s", serviceURL, r.URL.Path)

	proxy.ServeHTTP(w, r)
}

func discoverService(name string) (*url.URL, error) {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	services, _, err := consul.Health().Service(name, "", true, nil)
	if err != nil {
		return nil, fmt.Errorf("could not query Consul for service '%s': %w", name, err)
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("no healthy instances of service '%s' found in Consul", name)
	}

	service := services[0].Service
	serviceAddress := fmt.Sprintf("http://%s:%d", service.Address, service.Port)

	return url.Parse(serviceAddress)
}
