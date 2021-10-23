package drivers

import (
	"context"
	"errors"

	// "errors"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
)

func GetDriver(host steniciov1alpha1.SqlHost) (Driver, error) {
	switch host.Spec.Engine {
	case steniciov1alpha1.EngineTypeMysql:
		return &MySqlDriver{
			Host: host,
		}, nil
	}

	return nil, errors.New("Driver could not be resolved")
}

type Driver interface {
	UpsertUser(context.Context, steniciov1alpha1.SqlUser) error
	DeleteUser(context.Context, steniciov1alpha1.SqlUser) error
	UpsertDatabase(context.Context, steniciov1alpha1.SqlDatabase) error
	DeleteDatabase(context.Context, steniciov1alpha1.SqlDatabase) error
	UpsertGrants(context.Context, steniciov1alpha1.SqlGrant, steniciov1alpha1.SqlUser, steniciov1alpha1.SqlDatabase) error
	DeleteGrants(context.Context, steniciov1alpha1.SqlGrant, steniciov1alpha1.SqlUser, steniciov1alpha1.SqlDatabase) error
}
