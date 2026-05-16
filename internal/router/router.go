package router

import (
	"net/http"
	"time"

	"github.com/ChaitanyaSaiV/Incident-Management/internal/handlers"
)

func Routes(path string, i *handlers.IncidentHandler) *http.Server {

	r := http.NewServeMux()
	r.HandleFunc("GET /health", handlers.HealthCheck)
	r.HandleFunc("GET /incidents/{id}", i.GetIncident)
	r.HandleFunc("GET /incidents", i.GetAllIncidents)
	r.HandleFunc("POST /incidents", i.SaveIncident)
	r.HandleFunc("DELETE /incidents/{id}", i.DeleteIncident)

	return &http.Server{
		Addr:         path,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
