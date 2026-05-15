package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
)

func TestInMemorySave(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	incident := models.IncidentData{
		Id:        "1",
		Severity:  "SEV1",
		Message:   "Very Severe",
		TimeStamp: time.Now(),
	}
	store.Save(ctx, &incident)

	getIncident, err := store.Get(ctx, incident.Id)
	if err != nil {
		t.Fatal("unable to save the incident")
	}

	if getIncident.Severity != incident.Severity {
		t.Fatal("Issue saving the incident")
	}
}

func TestInMemoryGetAll(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	incident := models.IncidentData{
		Id:        "1",
		Severity:  "SEV1",
		Message:   "Very Severe",
		TimeStamp: time.Now(),
	}
	store.Save(ctx, &incident)
	incidents, err := store.GetAll(ctx)
	if err != nil {
		t.Fatal("error retrieving the incidents")
	}

	if len(incidents) != 1 {
		t.Fatal("error retrieving the incidents")
	}
}

func TestInMemoryConcurrentSaves(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	var wg sync.WaitGroup
	incidentCount := 100
	for i := 0; i < incidentCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			incident := models.IncidentData{
				Id:        fmt.Sprintf("%d", id),
				Severity:  "SEV1",
				Message:   fmt.Sprintf("Incident Number : %d", id),
				TimeStamp: time.Now(),
			}
			store.Save(ctx, &incident)
		}(i)
	}
	wg.Wait()

	incidents, _ := store.GetAll(ctx)
	if len(incidents) != incidentCount {
		t.Fatalf("issue with saving concurrent requests, expected %d got %d", incidentCount, len(incidents))
	}
}

func TestInMemoryConcurrentSave_ConcurrentRead(t *testing.T) {
	store := NewInMemoryStore()
	ctx := context.Background()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		incident := models.IncidentData{
			Id:        fmt.Sprintf("%d", i),
			Severity:  "SEV1",
			Message:   "Message",
			TimeStamp: time.Now(),
		}
		store.Save(ctx, &incident)
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			incident := models.IncidentData{
				Id:        fmt.Sprintf("%d", id),
				Severity:  "SEV1",
				Message:   "Message",
				TimeStamp: time.Now(),
			}
			store.Save(ctx, &incident)

		}(i)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		id := fmt.Sprintf("%d", i)
		go func(id string) {
			defer wg.Done()
			store.Get(ctx, id)

		}(id)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			store.GetAll(ctx)
		}()
	}

	wg.Wait()
}
