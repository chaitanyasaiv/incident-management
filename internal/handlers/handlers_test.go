package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
	"github.com/ChaitanyaSaiV/Incident-Management/internal/store"
)

type MockStore struct {
	// ─── INPUTS: configure what the mock returns ───
	GetData    models.IncidentData   // what Get should return on success
	GetErr     error                 // what Get should return as error (set to ErrNotFound to simulate missing)
	GetAllData []models.IncidentData // what GetAll should return
	GetAllErr  error                 // what GetAll should return as error
	SaveErr    error                 // what Save should return (nil = success)
	DeleteErr  error                 // what Delete should return

	// ─── OUTPUTS: observe what the mock was called with ───
	SavedIncident *models.IncidentData // captures the last incident passed to Save
	DeletedID     string               // captures the last ID passed to Delete
	GetCalledID   string               // captures the last ID passed to Get

	// ─── CALL COUNTS: verify how many times each method was called ───
	SaveCalled   int
	GetCalled    int
	GetAllCalled int
	DeleteCalled int
}

func (m *MockStore) Get(ctx context.Context, id string) (models.IncidentData, error) {
	m.GetCalled++
	m.GetCalledID = id
	return m.GetData, m.GetErr
}

func (m *MockStore) GetAll(ctx context.Context) ([]models.IncidentData, error) {
	m.GetAllCalled++
	return m.GetAllData, m.GetAllErr
}

func (m *MockStore) Save(ctx context.Context, incident *models.IncidentData) error {
	m.SaveCalled++
	m.SavedIncident = incident
	return m.SaveErr
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	m.DeleteCalled++
	m.DeletedID = id
	return m.DeleteErr
}

var _ store.IncidentStore = (*MockStore)(nil)

func TestSaveIncident(t *testing.T) {
	incidents := []struct {
		Name       string
		Data       string
		Err        error
		StatusCode int
	}{
		{
			Name:       "Save Incident",
			Data:       `{"id":"1","message":"db down","severity":"SEV1"}`,
			Err:        nil,
			StatusCode: http.StatusCreated,
		},
		{
			Name:       "Invalid Severity",
			Data:       `{"id":"1","message":"db down","severity":"SEV"}`,
			Err:        nil,
			StatusCode: http.StatusBadRequest,
		},
		{
			Name:       "Invalid id",
			Data:       `{"id":"1","message":"db down","severity":"SEV1"}`,
			Err:        nil,
			StatusCode: http.StatusCreated,
		},
	}

	for _, incident := range incidents {
		t.Run(incident.Name, func(t *testing.T) {
			mock := &MockStore{}
			handlers := NewHandler(mock)
			req := httptest.NewRequest(http.MethodPost, "/incidents",
				bytes.NewBufferString(incident.Data))
			rec := httptest.NewRecorder()
			handlers.SaveIncident(rec, req)
			if rec.Code != incident.StatusCode {
				t.Fatalf("Expected StatusCode : %v, actual statuscode : %v", incident.StatusCode, rec.Code)
			}
		})
	}
}
