package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/models"
)

type FileStorage struct {
	mu       sync.RWMutex
	fileName string
}

func NewFileStorage(fileName string) (*FileStorage, error) {
	_, err := os.Stat(fileName)
	if os.IsExist(err) {
		err = os.WriteFile(fileName, []byte{}, 0644)
		if err != nil {
			return nil, fmt.Errorf("Error creating the file %v", err)
		}
	}

	return &FileStorage{
		fileName: fileName,
	}, nil
}

func (fs *FileStorage) Get(ctx context.Context, id string) (models.IncidentData, error) {

	fs.mu.RLock()
	defer fs.mu.RUnlock()

	fileData, err := os.ReadFile(fs.fileName)
	if err != nil {
		return models.IncidentData{}, fmt.Errorf("error reading the file")
	}

	var incidents map[string]models.IncidentData

	json.Unmarshal(fileData, &incidents)

	incident, ok := incidents[id]
	if !ok {

		return models.IncidentData{}, fmt.Errorf("incident record not found")
	}

	return incident, nil

}

func (fs *FileStorage) GetAll(ctx context.Context) ([]models.IncidentData, error) {

	fs.mu.RLock()
	defer fs.mu.RUnlock()

	fileData, err := os.ReadFile(fs.fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading the file")
	}

	var incidents map[string]models.IncidentData

	json.Unmarshal(fileData, &incidents)

	var responseData []models.IncidentData

	for _, incident := range incidents {
		responseData = append(responseData, incident)
	}

	return responseData, nil

}
func (fs *FileStorage) Save(ctx context.Context, incident *models.IncidentData) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fileData, err := os.ReadFile(fs.fileName)
	if err != nil {
		return fmt.Errorf("error reading the file")
	}
	incidents := make(map[string]models.IncidentData)
	json.Unmarshal(fileData, &incidents)
	incidents[incident.Id] = *incident
	fileData, err = json.MarshalIndent(incidents, "", " ")
	if err != nil {
		return fmt.Errorf("internal server error")
	}
	err = os.WriteFile(fs.fileName, fileData, 0644)
	if err != nil {
		return fmt.Errorf("error writing the file")
	}
	return nil
}

func (fs *FileStorage) Delete(ctx context.Context, id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fileData, err := os.ReadFile(fs.fileName)
	if err != nil {
		return fmt.Errorf("error reading the file")
	}
	incidents := make(map[string]models.IncidentData)
	err = json.Unmarshal(fileData, &incidents)
	if err != nil {
		return fmt.Errorf("error un-marshalling the incidents data")
	}
	delete(incidents, id)

	newData, err := json.MarshalIndent(incidents, "", " ")

	if err != nil {
		return fmt.Errorf("error marshalling the incidents data")
	}

	err = os.WriteFile(fs.fileName, newData, 0644)

	if err != nil {
		return fmt.Errorf("error writing the file")
	}

	return nil

}
