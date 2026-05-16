package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/storage"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	responseData := models.HealthCheck{
		Status:    "OK",
		TimeStamp: time.Now().UTC(),
	}

	json.NewEncoder(w).Encode(responseData)
}

type IncidentStore interface {
	Get(ctx context.Context, id string) (models.IncidentData, error)
	GetAll(ctx context.Context) ([]models.IncidentData, error)
	Save(ctx context.Context, incident *models.IncidentData) error
	Delete(ctx context.Context, id string) error
}

type IncidentHandler struct {
	incidents IncidentStore
}

func NewHandler(s IncidentStore) *IncidentHandler {
	return &IncidentHandler{
		incidents: s,
	}
}

func (i *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	incident, err := i.incidents.Get(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(incident)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (i *IncidentHandler) SaveIncident(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateIncident

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	incident := models.IncidentData{
		Id:        reqBody.Id,
		Severity:  reqBody.Severity,
		Message:   reqBody.Message,
		TimeStamp: time.Now(),
	}

	err = i.incidents.Save(r.Context(), &incident)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("successfully saved the incident"))
}

func (i *IncidentHandler) GetAllIncidents(w http.ResponseWriter, r *http.Request) {
	allIncidents, err := i.incidents.GetAll(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(allIncidents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (i *IncidentHandler) DeleteIncident(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := i.incidents.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
