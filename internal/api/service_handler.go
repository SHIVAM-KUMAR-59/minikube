package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SHIVAM-KUMAR-59/minikube/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateServiceRequest struct {
	Name string   `json:"name"`
	Pods []string `json:"pods"`
	Port string   `json:"port"`
}

// CreateService handles the /services endpoint for creating a new Service. It decodes the request body to get the Service name, associated pods, and port, validates the input, creates a new Service struct, saves it to the store, and responds with a success message and the created Service.
func (h *Handler) CreateService(res http.ResponseWriter, req *http.Request) {
	var createServiceReq CreateServiceRequest
	err := json.NewDecoder(req.Body).Decode(&createServiceReq)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if createServiceReq.Name == "" || createServiceReq.Port == "" {
		slog.Error("Name and port are required fields")
		http.Error(res, "Name and port are required fields", http.StatusBadRequest)
		return
	}

	service := store.Service{
		ID:   generateRandomID(),
		Name: createServiceReq.Name,
		Pods: createServiceReq.Pods,
		Port: createServiceReq.Port,
	}

	err = h.store.CreateService(service)
	if err != nil {
		slog.Error("Failed to create service", "error", err)
		http.Error(res, "Failed to create service", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(service)
	slog.Info("Service created successfully", "service_id", service.ID)
}

// GetAllServices handles the /services endpoint for retrieving all Services. It retrieves all services from the store, encodes them as JSON, and responds with the list of services.
func (h *Handler) GetAllServices(res http.ResponseWriter, req *http.Request) {
	services, err := h.store.GetAllServices()
	if err != nil {
		slog.Error("Failed to retrieve services", "error", err)
		http.Error(res, "Failed to retrieve services", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	json.NewEncoder(res).Encode(services)
	slog.Info("Retrieved all services successfully", "service_count", len(services))
}

// GetNextPod handles the /services/{name}/next endpoint for retrieving the next pod for a given service. It extracts the service name from the URL, validates it, uses the load balancer to get the next pod for the service, and responds with the pod information in JSON format.
func (h *Handler) GetNextPod(res http.ResponseWriter, req *http.Request) {
	serviceName := chi.URLParam(req, "name")

	if serviceName == "" {
		slog.Error("Service name is required")
		http.Error(res, "Service name is required", http.StatusBadRequest)
		return
	}

	pod, err := h.loadBalancer.GetNextPodForService(serviceName)
	if err != nil {
		slog.Error("Failed to get next pod for service", "error", err)
		http.Error(res, "Failed to get next pod for service", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	json.NewEncoder(res).Encode(pod)
	slog.Info("Retrieved next pod for service successfully", "service_name", serviceName)
}

// DeleteService handles the /services/{id} endpoint for deleting a Service. It extracts the service ID from the URL, validates it, deletes the service from the store, and responds with a success message if the deletion is successful.
func (h *Handler) DeleteService(res http.ResponseWriter, req *http.Request) {
	serviceID := chi.URLParam(req, "id")

	if serviceID == "" {
		slog.Error("Service ID is required")
		http.Error(res, "Service ID is required", http.StatusBadRequest)
		return
	}

	err := h.store.DeleteService(serviceID)
	if err != nil {
		slog.Error("Failed to delete service", "error", err)
		http.Error(res, "Failed to delete service", http.StatusInternalServerError)
		return
	}

	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintln(res, `{"status": "ok", "message": "Service deleted successfully"}`)
	slog.Info("Service deleted successfully", "service_id", serviceID)
}
