package profile

import (
	"context"
	"rentCar/server/shared/id"
)

type Manager struct {
}

func (p *Manager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return id.IdentityID("identity1"), nil
}