package drivers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
)

type driverDecorator struct {
	driver        Driver
	Noop          bool
	ownerShipData OwnerShipData
}

func (d *driverDecorator) UpsertUser(ctx context.Context, user steniciov1alpha1.SqlUser) (int64, error) {
	if !d.Noop {
		return d.driver.UpsertUser(ctx, user)
	}
	return 0, nil
}

func (d *driverDecorator) DeleteUser(ctx context.Context, user steniciov1alpha1.SqlUser) error {
	if user.Spec.CleanupPolicy == steniciov1alpha1.CleanupPolicyDelete {
		return d.driver.DeleteUser(ctx, user)
	}
	return nil
}

func (d *driverDecorator) UpsertDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) (int64, error) {
	if !d.Noop {
		return d.driver.UpsertDatabase(ctx, database)
	}
	return 0, nil
}

func (d *driverDecorator) DeleteDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) error {
	if !d.Noop && database.Spec.CleanupPolicy == steniciov1alpha1.CleanupPolicyDelete {
		return d.driver.DeleteDatabase(ctx, database)
	}
	return nil
}

func (d *driverDecorator) UpsertGrants(ctx context.Context, grant steniciov1alpha1.SqlGrant, user steniciov1alpha1.SqlUser, database steniciov1alpha1.SqlDatabase) (int64, error) {
	if !d.Noop {
		return d.driver.UpsertGrants(ctx, grant, user, database)
	}
	return 0, nil
}

func (d *driverDecorator) DeleteGrants(ctx context.Context, grant steniciov1alpha1.SqlGrant, user steniciov1alpha1.SqlUser, database steniciov1alpha1.SqlDatabase) error {
	if !d.Noop && grant.Spec.CleanupPolicy == steniciov1alpha1.CleanupPolicyDelete {
		return d.driver.DeleteGrants(ctx, grant, user, database)
	}
	return nil
}

func (d *driverDecorator) DeleteOwnerState(ctx context.Context) error {
	return d.driver.DeleteOwnerState(ctx, d.ownerShipData)
}

func (d *driverDecorator) SetOwnershipData(ctx context.Context, data OwnerShipData) error {
	d.ownerShipData = data
	if d.ownerShipData.OwnerID == "" {
		d.ownerShipData.OwnerID = steniciov1alpha1.OwnerID(uuid.New().String())
	}

	err := d.checkOwner(ctx)
	if err != nil {
		d.Noop = true
	}

	return err
}

func (d *driverDecorator) InitOwnerSchema(ctx context.Context) error {
	return d.driver.InitOwnerSchema(ctx)
}

func (d *driverDecorator) GetOwnerID() steniciov1alpha1.OwnerID {
	return d.ownerShipData.OwnerID
}

func (d *driverDecorator) checkOwner(ctx context.Context) error {
	ownerState, err := d.driver.GetOwnerState(ctx, d.ownerShipData)
	if err != nil {
		return err
	}
	switch ownerState {
	case NonExisting:
		if err := d.driver.SetOwnerState(ctx, d.ownerShipData); err != nil {
			return fmt.Errorf("failed to claim ownership - %s", err.Error())
		}
		return nil
	case NotOwner:
		// We are not managing this resource)
		return fmt.Errorf("resource already exists and will not managed by this sqlOperator")
	case IsOwner:
		// All good, proceed
		return nil
	default:
		return fmt.Errorf("unknown ownerState '%s'", ownerState)
	}
}
