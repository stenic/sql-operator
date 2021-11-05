# Getting started

To see how Sql-operator works, you can install it and configure some examples.

Firstly, you'll need a Kubernetes cluster, kubectl and helm set-up

## Installing Sql-operator

To get started quickly, you can use the
[helm chart](https://artifacthub.io/packages/helm/sql-operator/sql-operator)
which will install Sql-operator:

```sh
helm repo add sql-operator https://stenic.github.io/sql-operator/
helm install sql-operator --namespace sql-operator --create-namespace sql-operator/sql-operator
```

!!! note
	The examples will use an in-cluster database for demo simplicity. Please use the following command
	to install it. 

	```sh
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm install test-db \
		--namespace sql-operator-testdb --create-namespace \
		--set auth.rootPassword=n0tS3cur3 \
		--set primary.persistence.enabled=false \
		--set secondary.persistence.enabled=false \
		bitnami/mariadb
	```

## Define an host you want to manage

First of, we will need to configure one or more `SqlHost` objects. These object are used by sql-operator
to grab the endpoint and the correct credentials.


```sh
kubectl apply -f https://raw.githubusercontent.com/stenic/sql-operator/master/examples/sqlhost-mysql.yaml
kubectl get sqlhost
```

!!! note
	If you didn't install the in-cluster database, you can alter the manifest.
	
	```sh
	kubectl edit sqlhost sample-host
	```

## Create a database, user and grants

Not that we have a host defined, we can create our first assets.

```sh
kubectl apply -f https://raw.githubusercontent.com/stenic/sql-operator/master/examples/sqldatabase-gettingstarted.yaml
kubectl apply -f https://raw.githubusercontent.com/stenic/sql-operator/master/examples/sqluser-gettingstarted.yaml
kubectl apply -f https://raw.githubusercontent.com/stenic/sql-operator/master/examples/sqlgrant-gettingstarted.yaml
```

Let's verify our resources got created. (Password is `n0tS3cur3`)

```sh
kubectl exec -ti -n sql-operator-testdb test-db-mariadb-0 -- mysql -uroot -p -e 'show databases;'
kubectl exec -ti -n sql-operator-testdb test-db-mariadb-0 -- mysql -uroot -p -e 'select user from mysql.user;'
kubectl exec -ti -n sql-operator-testdb test-db-mariadb-0 -- mysql -uroot -p -e 'show grants for sample_username;'
```
