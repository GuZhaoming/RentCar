package car

import (
	"context"
	rentalpb "rentCar/server/rental/api/gen/v1"
	"rentCar/server/shared/id"
)

type Manager struct {
}

func (c *Manager) Verify(context.Context, id.CarID, *rentalpb.Location) error {
	return nil
}

func (c *Manager) Unlock(context.Context, id.CarID) error {
	return nil
}