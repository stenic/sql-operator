# Core concepts

Sql-operator aims to solve 3 different cases when dealing with Sql resources.

## Create

Simplify creation of Sql resouces that are dependecies of your application. 
You define what database, user and grants you want and Sql-operator will ensure they are created.

## Reconcile

At set intervals (default 10 seconds) sql-operator will validate if the resources are still in the
desired state. If not it will re-apply the wanted state. This adds an additional layer of control.

## Cleanup

If configured using the `cleanupPolicy`, sql-operator will delete the external resouces. This way
your resources stay clean.
