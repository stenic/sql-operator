package controllers

import (
	"context"
	"fmt"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"github.com/stenic/sql-operator/drivers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func checkOwner(ctx context.Context, driver drivers.Driver, ownerShipData drivers.OwnerShipData, obj client.Object) (steniciov1alpha1.OwnerID, error) {
	ownerState, err := driver.GetOwnerState(ctx, ownerShipData)
	switch ownerState {
	case drivers.NonExisting:
		ownerShipData.OwnerID = steniciov1alpha1.OwnerID(obj.GetUID())
		if err := driver.SetOwnerState(ctx, ownerShipData); err != nil {
			return "", fmt.Errorf("failed to claim ownership - %s", err.Error())
		}
		return ownerShipData.OwnerID, nil
	case drivers.NotOwner:
		// We are not managing this resource)
		return "", fmt.Errorf("resource already exists and will not managed by this sqlOperator")
	case drivers.IsOwner:
		// All good, proceed
		break
	default:
		return "", fmt.Errorf("unknown ownerState '%s'", ownerState)
	}

	return "", err
}

func cleanupOwner(ctx context.Context, driver drivers.Driver, ownerShipData drivers.OwnerShipData) error {
	return driver.DeleteOwnerState(ctx, ownerShipData)
}
