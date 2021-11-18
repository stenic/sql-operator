package drivers

import (
	"context"
	"errors"

	// "errors"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
)

func GetDriver(host steniciov1alpha1.SqlHost) (*driverDecorator, error) {
	var driver Driver
	switch host.Spec.Engine {
	case steniciov1alpha1.EngineTypeMysql:
		driver = &MySqlDriver{
			Host: host,
		}
	}

	if driver != nil {
		return &driverDecorator{
			driver: driver,
		}, nil
	}

	return nil, errors.New("Driver could not be resolved")
}

type OwnerShipType string

type OwnerShipData struct {
	Type     OwnerShipType
	Name     string
	Resource string
	OwnerID  steniciov1alpha1.OwnerID
}

type OwnerState string

const (
	IsOwner     OwnerState = "IsOwner"
	NotOwner    OwnerState = "NotOwner"
	NonExisting OwnerState = "NonExisting"

	OwnerShipTypeDatabase OwnerShipType = "OwnerShipTypeDatabase"
	OwnerShipTypeUser     OwnerShipType = "OwnerShipTypeUser"
	OwnerShipTypeGrant    OwnerShipType = "OwnerShipTypeGrant"
)

type Driver interface {
	UpsertUser(context.Context, steniciov1alpha1.SqlUser) (int64, error)
	DeleteUser(context.Context, steniciov1alpha1.SqlUser) error
	UpsertDatabase(context.Context, steniciov1alpha1.SqlDatabase) (int64, error)
	DeleteDatabase(context.Context, steniciov1alpha1.SqlDatabase) error
	UpsertGrants(context.Context, steniciov1alpha1.SqlGrant, steniciov1alpha1.SqlUser, steniciov1alpha1.SqlDatabase) (int64, error)
	DeleteGrants(context.Context, steniciov1alpha1.SqlGrant, steniciov1alpha1.SqlUser, steniciov1alpha1.SqlDatabase) error
	InitOwnerSchema(context.Context) error
	SetOwnerState(context.Context, OwnerShipData) error
	DeleteOwnerState(context.Context, OwnerShipData) error
	GetOwnerState(context.Context, OwnerShipData) (OwnerState, error)
}
