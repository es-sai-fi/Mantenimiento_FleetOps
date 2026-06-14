package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/fleetops/maintenance/internal/domain"
	"github.com/fleetops/maintenance/internal/port"
)

// Compile-time check: HTTPVehicleClient implements port.VehicleClient.
var _ port.VehicleClient = (*HTTPVehicleClient)(nil)

// HTTPVehicleClient is the concrete implementation of port.VehicleClient
// that communicates with the external Vehicles microservice via REST/HTTP.
//
// This adapter implements the Anti-Corruption Layer (ACL), translating
// external API responses into domain.Vehicle value objects.
//
// [Archetype Convention Addition] — Anti-Corruption Layer (DDD best practice)
// SAD Reference: Process Network 2 — "GET /vehiculos", "PUT /vehiculos"
type HTTPVehicleClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPVehicleClient constructs an HTTPVehicleClient.
func NewHTTPVehicleClient(baseURL string, timeoutSecs int) *HTTPVehicleClient {
	return &HTTPVehicleClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSecs) * time.Second,
		},
	}
}

// externalVehicle represents the JSON structure returned by the external
// Vehicles microservice. This is the "foreign" model that gets translated.
type externalVehicle struct {
	ID                       string  `json:"id"`
	LicensePlate             string  `json:"license_plate"`
	KilometersRecorded       float64 `json:"kilometers_recorded"`
	DaysSinceLastMaintenance int     `json:"days_since_last_maintenance"`
	Available                bool    `json:"available"`
}

// GetAllVehicles fetches all vehicles from the Vehicles microservice and
// translates them into domain.Vehicle value objects.
//
// ACL Translation: externalVehicle → domain.Vehicle
func (c *HTTPVehicleClient) GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating vehicle list request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing vehicle list request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vehicle service returned status %d", resp.StatusCode)
	}

	var external []externalVehicle
	if err := json.NewDecoder(resp.Body).Decode(&external); err != nil {
		return nil, fmt.Errorf("decoding vehicle list response: %w", err)
	}

	// ACL Translation: convert external models to domain models
	vehicles := make([]*domain.Vehicle, 0, len(external))
	for _, ev := range external {
		id, err := uuid.Parse(ev.ID)
		if err != nil {
			continue // skip vehicles with invalid IDs
		}
		vehicles = append(vehicles, &domain.Vehicle{
			ID:                       id,
			LicensePlate:             ev.LicensePlate,
			KilometersRecorded:       ev.KilometersRecorded,
			DaysSinceLastMaintenance: ev.DaysSinceLastMaintenance,
			Available:                ev.Available,
		})
	}

	return vehicles, nil
}

// vehicleUpdatePayload represents the JSON body sent to the Vehicles
// microservice when updating maintenance status.
type vehicleUpdatePayload struct {
	DaysSinceLastMaintenance int `json:"days_since_last_maintenance"`
}

// UpdateVehicleMaintenanceStatus sends a PUT request to the Vehicles
// microservice to reset the maintenance days counter after a maintenance
// is completed.
//
// SAD Reference: "PUT a /vehículos — actualiza estado y cantidad de días
// desde el último mantenimiento del vehículo"
func (c *HTTPVehicleClient) UpdateVehicleMaintenanceStatus(ctx context.Context, vehicleID uuid.UUID, daysReset int) error {
	payload := vehicleUpdatePayload{DaysSinceLastMaintenance: daysReset}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshaling vehicle update payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", c.baseURL, vehicleID.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating vehicle update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing vehicle update request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("vehicle service returned status %d on update", resp.StatusCode)
	}

	return nil
}
