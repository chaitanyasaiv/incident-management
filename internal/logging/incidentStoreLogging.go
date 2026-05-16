package logging

import (
	"context"
	"log"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/handlers"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
)

type LoggingIncidentStore struct {
	next handlers.IncidentStore
}

func NewLoggingIncidentStore(store handlers.IncidentStore) *LoggingIncidentStore {
	return &LoggingIncidentStore{
		next: store,
	}
}

func (l *LoggingIncidentStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	start := time.Now()
	incident, err := l.next.Get(ctx, id)
	log.Printf("Get id=%s err=%v took=%v", id, err, time.Since(start))
	return incident, err
}
func (l *LoggingIncidentStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	start := time.Now()
	incidents, err := l.next.GetAll(ctx)
	log.Println(time.Since(start))
	return incidents, err
}
func (l *LoggingIncidentStore) Save(ctx context.Context, incident *models.IncidentData) error {
	start := time.Now()
	err := l.next.Save(ctx, incident)
	log.Println(time.Since(start))
	return err
}
func (l *LoggingIncidentStore) Delete(ctx context.Context, id string) error {
	start := time.Now()
	err := l.next.Delete(ctx, id)
	log.Println(time.Since(start))
	return err
}
