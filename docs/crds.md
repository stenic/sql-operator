# CRDs


```yaml
# sqlhosts.stenic.io
apiVersion: stenic.io/v1alpha1
kind: SqlHost
metadata:
  name: integration
spec:
  engine: Mysql
  dsn: sqloperator:xxxxxxxxxxx@tcp(192.168.1.123:3306)/
```

```yaml
# sqldatabase.stenic.io
apiVersion: stenic.io/v1alpha1
kind: SqlDatabase
metadata:
  name: application-int-db
spec:
  databaseName: myapp_testing
  hostRef:
    name: integration
  cleanupPolicy: Delete
```

```yaml
# sqluser.stenic.io
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
```


```yaml
# sqlgrant.stenic.io
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
