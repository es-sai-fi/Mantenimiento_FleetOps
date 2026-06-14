package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles the /health endpoint for service health checks.
//
// [Archetype Convention Addition] — Health check endpoint for Docker
// HEALTHCHECK and load balancer probes (AWS ECS best practice).
type HealthHandler struct {
	pool *pgxpool.Pool
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

// healthResponse represents the health check JSON response.
type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

// Check handles GET /health
// Returns service status and database connectivity.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	dbStatus := "up"

	if err := h.pool.Ping(r.Context()); err != nil {
		dbStatus = "down"
	}

	status := "healthy"
	statusCode := http.StatusOK
	if dbStatus == "down" {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:   status,
		Database: dbStatus,
	})
}
