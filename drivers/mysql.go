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

func (d *MySqlDriver) connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/",
		d.Host.Spec.Credentials.Username,
		d.Host.Spec.Credentials.Password,
		d.Host.Spec.Endpoint.Host,
		d.Host.Spec.Endpoint.Port,
	)
	db, err := sql.Open("mysql", dsn)
	if err == nil {
		db.SetConnMaxLifetime(time.Minute * 1)
	}
	return db, err
}

func (d *MySqlDriver) UpsertUser(ctx context.Context, user steniciov1alpha1.SqlUser) error {
	log := log.FromContext(ctx)
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Try to alter user first
	if _, err := db.ExecContext(ctx, fmt.Sprintf(
		"ALTER USER '%s'@'%%' IDENTIFIED BY '%s';",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	)); err != nil {
		// User may not exist, try to create it
		if _, err := db.ExecContext(ctx, fmt.Sprintf(
			"CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s';",
			user.Spec.Credentials.Username,
			user.Spec.Credentials.Password,
		)); err != nil {
			return err
		}
	}

	log.V(1).Info("UPSERT USER")

	return nil
}

func (d *MySqlDriver) DeleteUser(ctx context.Context, user steniciov1alpha1.SqlUser) error {
	log := log.FromContext(ctx)

	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the user
	if _, err := db.ExecContext(ctx, fmt.Sprintf(
		"DROP USER IF EXISTS '%s'@'%%';",
		user.Spec.Credentials.Username,
	)); err != nil {
		return err
	}

	log.V(1).Info("DELETE USER")

	return nil
}

func (d *MySqlDriver) UpsertDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) error {
	log := log.FromContext(ctx)
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// User may not exist, try to create it
	if _, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s;",
		database.Spec.DatabaseName,
	)); err != nil {
		return err
	}

	log.V(1).Info("UPSERT DATABASE")

	return nil
}

func (d *MySqlDriver) DeleteDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) error {
	log := log.FromContext(ctx)

	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the database
	if _, err := db.QueryContext(ctx, fmt.Sprintf(
		"DROP DATABASE IF EXISTS %s;",
		database.Spec.DatabaseName,
	)); err != nil {
		return err
	}

	log.V(1).Info("DELETE DATABASE")

	return nil
}
