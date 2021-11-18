package drivers

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var refSchema = "sqloperator_schema"
var refSchemaDatabase = "sqldatabase_ref"
var refSchemaUser = "sqluser_ref"
var refSchemaGrant = "sqlgrant_ref"

var createOwnerSchema = fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s;`, refSchema)
var createOwnerSchemaTableTpl = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s.%%s (
  id varchar(64) NOT NULL COMMENT 'Unique identifier',
  name varchar(253) NOT NULL COMMENT 'Sql object name',
  owner varchar(253) NOT NULL COMMENT 'Controller name',
  resource varchar(507) NOT NULL COMMENT 'Kubernetes object reference',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`, refSchema)

type MySqlDriver struct {
	Host steniciov1alpha1.SqlHost
}

func (d *MySqlDriver) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", d.Host.Spec.DSN)
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

func (d *MySqlDriver) InitOwnerSchema(ctx context.Context) error {
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err = db.ExecContext(ctx, createOwnerSchema); err != nil {
		return err
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(createOwnerSchemaTableTpl, refSchemaDatabase)); err != nil {
		return err
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(createOwnerSchemaTableTpl, refSchemaUser)); err != nil {
		return err
	}
	if _, err = db.ExecContext(ctx, fmt.Sprintf(createOwnerSchemaTableTpl, refSchemaGrant)); err != nil {
		return err
	}
	return nil
}

func (d *MySqlDriver) SetOwnerState(ctx context.Context, data OwnerShipData) error {
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(
		ctx,
		fmt.Sprintf("INSERT IGNORE INTO %s (id, resource, name, owner) VALUES (?, ?, ?, ?) ;", d.getOwnerDBTable(data)),
		data.OwnerID, data.Resource, data.Name, "sql-operator",
	)

	return err
}

func (d *MySqlDriver) GetOwnerState(ctx context.Context, data OwnerShipData) (OwnerState, error) {
	state := NonExisting
	db, err := d.connect()
	if err != nil {
		return state, err
	}
	defer db.Close()

	sqlStatement := fmt.Sprintf("SELECT s.id FROM %s s WHERE resource=? LIMIT 1;", d.getOwnerDBTable(data))
	var dbID string
	row := db.QueryRowContext(ctx, sqlStatement, data.Resource)
	switch err := row.Scan(&dbID); err {
	case nil:
		if dbID == string(data.OwnerID) {
			return IsOwner, nil
		}
		return NotOwner, nil
	case sql.ErrNoRows:
		break
	default:
		return NotOwner, err
	}

	switch data.Type {
	case OwnerShipTypeDatabase:
		results, err := db.QueryContext(ctx, "SHOW databases;")
		if err != nil {
			return state, err
		}
		defer results.Close()

		for results.Next() {
			var dbName string
			if err = results.Scan(&dbName); err != nil {
				return state, err
			}
			if dbName == data.Name {
				return NotOwner, nil
			}
		}
	case OwnerShipTypeUser:
		sqlStatement := fmt.Sprintf("SELECT count(u.id) FROM %s s WHERE resource=? LIMIT 1;", d.getOwnerDBTable(data))
		var matchCount int
		row := db.QueryRowContext(ctx, sqlStatement, data.Resource)
		if err := row.Scan(&matchCount); err != nil {
			return NotOwner, err
		}
		if matchCount >= 1 {
			return NotOwner, nil
		}
	}

	return NonExisting, nil
}

func (d *MySqlDriver) DeleteOwnerState(ctx context.Context, data OwnerShipData) error {
	db, err := d.connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(
		ctx,
		fmt.Sprintf("DELETE FROM %s WHERE id = ? AND resource = ? LIMIT 1;", d.getOwnerDBTable(data)),
		data.OwnerID, data.Resource,
	)

	return err
}

func (d *MySqlDriver) getOwnerDBTable(data OwnerShipData) string {
	switch data.Type {
	case OwnerShipTypeDatabase:
		return fmt.Sprintf("`%s`.%s", refSchema, refSchemaDatabase)
	case OwnerShipTypeUser:
		return fmt.Sprintf("`%s`.%s", refSchema, refSchemaUser)
	case OwnerShipTypeGrant:
		return fmt.Sprintf("`%s`.%s", refSchema, refSchemaGrant)
	default:
		panic("unhandled ownership type")
	}
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

	for _, grantName := range grants.Spec.Grants {
		if _, err := db.ExecContext(ctx, fmt.Sprintf(
			"REVOKE %s ON %s.* FROM '%s'@'%%';",
			grantName,
			database.Spec.DatabaseName,
			user.Spec.Credentials.Username,
		)); err != nil {
			log.Error(err, "Failed to revoke "+grantName)
		}
	}

	log.V(1).Info("DELETE GRANTS")

	return nil
}
