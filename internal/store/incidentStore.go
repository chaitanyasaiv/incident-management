package store

import (
	"context"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
)

type IncidentStore interface {
	IncidentReader
	IncidentWriter
}

type IncidentReader interface {
	Get(ctx context.Context, id string) (models.IncidentData, error)
	GetAll(ctx context.Context) ([]models.IncidentData, error)
}

type IncidentWriter interface {
	Save(ctx context.Context, incident *models.IncidentData) error
	Delete(ctx context.Context, id string) error
}
