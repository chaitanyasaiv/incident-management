package models

import "time"

type CreateIncident struct {
	Id       string `json:"id" validate:"required"`
	Severity string `json:"severity" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

type IncidentData struct {
	Id        string `json:"id"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	TimeStamp time.Time
}
