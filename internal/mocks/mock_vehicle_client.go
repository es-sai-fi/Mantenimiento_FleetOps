package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/fleetops/maintenance/internal/domain"
)

// MockVehicleClient is a testify mock implementing port.VehicleClient.
// Generated manually to match mockery output format.
type MockVehicleClient struct {
	mock.Mock
}

func (m *MockVehicleClient) GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Vehicle), args.Error(1)
}

func (m *MockVehicleClient) UpdateVehicleMaintenanceStatus(ctx context.Context, vehicleID uuid.UUID, daysReset int) error {
	args := m.Called(ctx, vehicleID, daysReset)
	return args.Error(0)
}
