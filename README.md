# SqlOperator

[![Last release](https://github.com/stenic/sql-operator/actions/workflows/release.yaml/badge.svg)](https://github.com/stenic/sql-operator/actions/workflows/release.yaml)
[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/sql-operator)](https://artifacthub.io/packages/search?repo=sql-operator)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fstenic%2Fsql-operator.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fstenic%2Fsql-operator?ref=badge_shield)


Operate sql databases, users and grants.

> This is a WIP project and should not at all be used in production at this time.
> Feel free to validate in CI pipelines.

## Engines

Currenty, only MySQL is supported.

## Installation

```shell
kubectl create namespace sql-operator
helm repo add sql-operator https://stenic.github.io/sql-operator/
helm install sql-operator --namespace sql-operator sql-operator/sql-operator
```

## Objects

### SqlHost

The `SqlHost` object contains information how the operator should connect to the remote server. 
Note that this information should be protected using rbac.

### SqlDatabase

The `SqlDatabase` object manages a database on the referenced `SqlHost`.

### SqlUser

The `SqlUser` object manages a user on the referenced `SqlHost`.

### SqlGrant

The `SqlGrant` object manages the grant for the referenced `SqlUser` on the referenced `SqlDatabase`.

## Example

The following example registeres a `SqlHost` pointing to a shared integration host and creates a `database`,
`user` and `grants` to run integration tests for the application.

```yaml
---
# Register a host - Not created by the operator
apiVersion: stenic.io/v1alpha1
kind: SqlHost
metadata:
  name: integration
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
  name: application-int-db
spec:
  databaseName: myapp_testing
  hostRef:
    name: integration
  cleanupPolicy: Delete

---
# Create a user on the host
apiVersion: stenic.io/v1alpha1
kind: SqlUser
metadata:
  name: application-int-user
spec:
  credentials:
    username: myapp_username
    password: myapp_password
  hostRef:
    name: integration
  cleanupPolicy: Delete

---
# Add some grant to the user
apiVersion: stenic.io/v1alpha1
kind: SqlGrant
metadata:
  name: application-int-grant
spec:
  userRef:
    name: application-int-user
  databaseRef:
    name: application-int-db
  grants:
    - INSERT
    - SELECT
    - UPDATE
    - DELETE
  cleanupPolicy: Delete
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fstenic%2Fsql-operator.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fstenic%2Fsql-operator?ref=badge_large)