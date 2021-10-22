package drivers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type MySqlDriver struct {
	Host steniciov1alpha1.SqlHost
}

func (d MySqlDriver) UpsertUser(ctx context.Context, user steniciov1alpha1.SqlUser) error {
	log := log.FromContext(ctx)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/",
		d.Host.Spec.Credentials.Username,
		d.Host.Spec.Credentials.Password,
		d.Host.Spec.Endpoint.Host,
		d.Host.Spec.Endpoint.Port,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 1)

	// Try to alter user first
	if _, err := db.QueryContext(ctx, fmt.Sprintf(
		"ALTER USER '%s'@'%%' IDENTIFIED BY '%s';",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	)); err != nil {
		// User may not exist, try to create it
		if _, err := db.QueryContext(ctx, fmt.Sprintf(
			"CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s';",
			user.Spec.Credentials.Username,
			user.Spec.Credentials.Password,
		)); err != nil {
			return err
		}
	}

	log.V(1).Info("Upsert from driver")

	return nil
}
