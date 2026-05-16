package service

import (
	"context"
	"fmt"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/store"
)

type IncidentServices interface {
	GetIncident(ctx context.Context, id string) (models.IncidentData, error)
	CreateIncident(ctx context.Context, incident *models.IncidentData) error
	ListIncidents(ctx context.Context) ([]models.IncidentData, error)
	DeleteIncident(ctx context.Context, id string) error
}

type IncidentService struct {
	incidents store.IncidentStore
}

func NewIncidentService(store store.IncidentStore) *IncidentService {
	return &IncidentService{
		incidents: store,
	}
}

func (i *IncidentService) GetIncident(ctx context.Context, id string) (models.IncidentData, error) {
	incident, err := i.incidents.Get(ctx, id)
	if err != nil {
		return models.IncidentData{}, fmt.Errorf("Error while retrieving the incident")
	}
	return incident, nil
}

func (i *IncidentService) CreateIncident(ctx context.Context, incident *models.IncidentData) error {
	err := i.incidents.Save(ctx, incident)
	if err != nil {
		return fmt.Errorf("error saving the incident")
	}
	return nil
}

func (i *IncidentService) ListIncidents(ctx context.Context) ([]models.IncidentData, error) {
	incidents, err := i.incidents.GetAll(ctx)
	if err != nil {
		return []models.IncidentData{}, fmt.Errorf("error retrieving the incidents")
	}
	return incidents, nil
}

func (i *IncidentService) DeleteIncident(ctx context.Context, id string) error {
	err := i.incidents.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting the incident")
	}
	return nil
}
