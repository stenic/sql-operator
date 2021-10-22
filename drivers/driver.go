package drivers

import (
	"context"
	// "errors"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
)

func GetDriver(host steniciov1alpha1.SqlHost) (Driver, error) {
	return MySqlDriver{
		Host: host,
	}, nil
	// return nil, errors.New("Driver could not be resolved")
}

type Driver interface {
	UpsertUser(ctx context.Context, user steniciov1alpha1.SqlUser) error
}
