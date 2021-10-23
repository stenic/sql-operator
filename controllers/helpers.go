package controllers

import (
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
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
