# Ownership

Sql-operator comes with ownership tracking since
[v1.13.0](https://github.com/stenic/sql-operator/releases/tag/v1.13.0). Ownership tracking ensures that only
1 object can alter the Sql resources.

Each time an object gets reconciled, a check will be performed to validate that the kubernetes object is indeed
the object managing the Sql resource. An additional status field has been added called `OwnerID`. The tracking 
is implemented on each driver.

The following checks are performed:

```
Is the object known in the reference table?
  -> YES: Does the OwnerID match?
    -> YES: IsOwner
	-> NO: NotOwner
  -> NO: Does the resource exist?
    -> YES: NotOwner
	-> NO: NonExisting
```
