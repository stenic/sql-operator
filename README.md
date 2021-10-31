# SqlOperator

[![Last release](https://github.com/stenic/sql-operator/actions/workflows/release.yaml/badge.svg)](https://github.com/stenic/sql-operator/actions/workflows/release.yaml)


Operate sql databases, users and grants.

> This is a WIP project and should not at all be used in production at this time.
> Feel free to validate in CI pipelines.

# Engines

Currenty, only MySQL is supported.

## Example

```yaml
---
# Register a host - Not created by the operator
apiVersion: stenic.io/v1alpha1
kind: SqlHost
metadata:
  name: sqlhost-sample
spec:
  engine: Mysql
  endpoint:
    host: 192.168.1.123
    port: 3306
  credentials:
    username: sqloperator
    password: xxxxxxxxxxx

---
# Create a database on the host
apiVersion: stenic.io/v1alpha1
kind: SqlDatabase
metadata:
  name: sqldatabase-sample
spec:
  databaseName: test123
  hostRef:
    name: sqlhost-sample
  cleanupPolicy: Delete

---
# Create a user on the host
apiVersion: stenic.io/v1alpha1
kind: SqlUser
metadata:
  name: sqluser-sample
spec:
  credentials:
    username: sqloperator_tst
    password: sqloperator_tst
  hostRef:
    name: sqlhost-sample
  cleanupPolicy: Delete

---
# Add some grant to the user
apiVersion: stenic.io/v1alpha1
kind: SqlGrant
metadata:
  name: sqlgrants-sample
spec:
  userRef:
    name: sqluser-sample
  databaseRef:
    name: sqldatabase-sample
  grants:
    - INSERT
    - SELECT
    - UPDATE
    - DELETE
  cleanupPolicy: Delete
```

