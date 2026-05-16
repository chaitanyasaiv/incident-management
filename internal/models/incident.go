package models

import (
	"fmt"
	"time"
)

type CreateIncident struct {
	Id       string `json:"id" validate:"required"`
	Severity string `json:"severity" validate:"required,oneof=SEV1 SEV2 SEV3"`
	Message  string `json:"message" validate:"required"`
}

type IncidentData struct {
	Id        string    `json:"id" validate:"required"`
	Severity  string    `json:"severity" validate:"required,oneof=SEV1 SEV2 SEV3"`
	Message   string    `json:"message" validate:"required"`
	TimeStamp time.Time `validate:"required"`
}

func (c *CreateIncident) Validate() error {
	if c.Severity != "SEV1" && c.Severity != "SEV2" && c.Severity != "SEV3" {
		return fmt.Errorf("Severity must of value SEV1, SEV2 or SEV3")
	}
	return nil
}
