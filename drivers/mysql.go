package drivers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(5)
	}
	return db, err
}

func (d *MySqlDriver) UpsertUser(ctx context.Context, user steniciov1alpha1.SqlUser) (int64, error) {
	log := log.FromContext(ctx)
	db, err := d.connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Try to alter user first
	res, err := db.ExecContext(ctx, fmt.Sprintf(
		"ALTER USER '%s'@'%%' IDENTIFIED BY '%s';",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	))
	if err != nil {
		// User may not exist, try to create it
		res, err = db.ExecContext(ctx, fmt.Sprintf(
			"CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s';",
			user.Spec.Credentials.Username,
			user.Spec.Credentials.Password,
		))
		if err != nil {
			return 0, err
		}
	}

	rowCount, _ := res.RowsAffected()
	log.V(1).Info("UPSERT USER")

	return rowCount, nil
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

func (d *MySqlDriver) UpsertDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) (int64, error) {
	log := log.FromContext(ctx)
	db, err := d.connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Create database
	res, err := db.ExecContext(ctx, fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s;",
		database.Spec.DatabaseName,
	))
	if err != nil {
		return 0, err
	}

	rowsCount, _ := res.RowsAffected()
	log.V(1).Info("UPSERT DATABASE")

	return rowsCount, nil
}

func (d *MySqlDriver) DeleteDatabase(ctx context.Context, database steniciov1alpha1.SqlDatabase) error {
	log := log.FromContext(ctx)

	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Delete the database
	if _, err := db.ExecContext(ctx, fmt.Sprintf(
		"DROP DATABASE IF EXISTS %s;",
		database.Spec.DatabaseName,
	)); err != nil {
		return err
	}

	log.V(1).Info("DELETE DATABASE")

	return nil
}

func (d *MySqlDriver) UpsertGrants(ctx context.Context, grants steniciov1alpha1.SqlGrant, user steniciov1alpha1.SqlUser, database steniciov1alpha1.SqlDatabase) (int64, error) {
	log := log.FromContext(ctx)
	db, err := d.connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	grantsQuery := fmt.Sprintf("SHOW GRANTS FOR '%s'@'%%';", user.Spec.Credentials.Username)
	results, err := db.QueryContext(ctx, grantsQuery)
	if err != nil {
		return 0, err
	}
	defer results.Close()

	r := regexp.MustCompile(fmt.Sprintf(
		`^GRANT\s(.+)\sON .?%s.?`,
		database.Spec.DatabaseName,
	))

	var currentGrants []string
	for results.Next() {
		var data string
		if err = results.Scan(&data); err != nil {
			return 0, err
		}
		if r.MatchString(data) {
			currentGrants = append(currentGrants, strings.Split(r.FindStringSubmatch(data)[1], ", ")...)
		}
	}

	toGrant := difference(grants.Spec.Grants, currentGrants)
	toRevoke := difference(currentGrants, grants.Spec.Grants)

	log.V(1).Info(fmt.Sprintf("toGrant:  %v\n", pp(toGrant)))
	log.V(1).Info(fmt.Sprintf("toRevoke: %v\n", pp(toRevoke)))

	// Revoke removed
	for _, grantName := range toRevoke {
		if _, err := db.ExecContext(ctx, fmt.Sprintf(
			"REVOKE %s ON %s.* FROM '%s'@'%%';",
			grantName,
			database.Spec.DatabaseName,
			user.Spec.Credentials.Username,
		)); err != nil {
			log.Error(err, "Failed to revoke "+grantName)
		}
	}

	// Grant added
	for _, grantName := range toGrant {
		if _, err := db.ExecContext(ctx, fmt.Sprintf(
			"GRANT %s ON %s.* TO '%s'@'%%';",
			grantName,
			database.Spec.DatabaseName,
			user.Spec.Credentials.Username,
		)); err != nil {
			log.Error(err, "Failed to grant "+grantName)
		}
	}

	changeCount := int64(len(toGrant) + len(toRevoke))
	// Flush when something changed
	if changeCount > 0 {
		if _, err := db.ExecContext(ctx, "FLUSH PRIVILEGES"); err != nil {
			return changeCount, fmt.Errorf("%v: %v", "Failed to flush", err)
		}
	}

	log.V(1).Info("UPSERT GRANTS")

	// err may have been set by a failed grant/revoke call caused by missing permission on the SqlHost.
	return changeCount, err
}

func pp(a []string) string {
	j, _ := json.Marshal(a)
	return string(j)
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func (d *MySqlDriver) DeleteGrants(ctx context.Context, grants steniciov1alpha1.SqlGrant, user steniciov1alpha1.SqlUser, database steniciov1alpha1.SqlDatabase) error {
	log := log.FromContext(ctx)

	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	// // Delete the database
	// if _, err := db.QueryContext(ctx, fmt.Sprintf(
	// 	"DROP DATABASE IF EXISTS %s;",
	// 	database.Spec.DatabaseName,
	// )); err != nil {
	// 	return err
	// }

	log.V(1).Info("DELETE GRANTS")

	return nil
}
