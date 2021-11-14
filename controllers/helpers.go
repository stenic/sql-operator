package controllers

import (
	"context"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	referencedDatabaseKey = ".spec.databaseRef"
	referencedUserKey     = ".spec.userRef"
	referencedHostKey     = ".spec.hostRef"
)

func getNamespacedName(ref steniciov1alpha1.SqlObjectRef, ns string) types.NamespacedName {
	if ref.Namespace != "" {
		return types.NamespacedName{
			Namespace: ref.Namespace,
			Name:      ref.Name,
		}
	}
	return types.NamespacedName{
		Namespace: ns,
		Name:      ref.Name,
	}
}

func isReferenced(ctx context.Context, c client.Client, children client.ObjectList, indexName string, obj client.Object) error {
	if err := c.List(ctx, children, client.MatchingFields{
		indexName: obj.GetNamespace() + "/" + obj.GetName(),
	}); err != nil {
		return err
	}

	return nil
}
