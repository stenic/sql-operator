/*
Copyright 2021 Stenic BV.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"github.com/stenic/sql-operator/drivers"
)

// SqlGrantReconciler reconciles a SqlGrant object
type SqlGrantReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Recorder    record.EventRecorder
	RefreshRate time.Duration
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqlgrants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqlgrants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqlgrants/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SqlGrantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var grants steniciov1alpha1.SqlGrant
	if err := r.Get(ctx, req.NamespacedName, &grants); err != nil {
		// log.Error(err, "unable to fetch SqlGrant")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	scheduledResult := ctrl.Result{RequeueAfter: r.RefreshRate}

	var user steniciov1alpha1.SqlUser
	if err := r.Get(ctx, getNamespacedName(grants.Spec.UserRef, grants.Namespace), &user); err != nil {
		log.Error(err, "unable to find SqlUser for "+grants.Name)
		return scheduledResult, client.IgnoreNotFound(err)
	}
	var database steniciov1alpha1.SqlDatabase
	if err := r.Get(ctx, getNamespacedName(grants.Spec.DatabaseRef, grants.Namespace), &database); err != nil {
		log.Error(err, "unable to find SqlDatabase for "+grants.Name)
		r.Recorder.Event(&grants, "Warning", "Error", "unable to find SqlDatabase")
		return scheduledResult, client.IgnoreNotFound(err)
	}
	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, getNamespacedName(database.Spec.HostRef, grants.Namespace), &host); err != nil {
		log.Error(err, "unable to find SqlHost for "+grants.Name)
		log.Info(fmt.Sprintf("DBG: Trying to find SqlHost with %v and %v ", database.Spec.HostRef, database.Spec.HostRef))
		r.Recorder.Event(&grants, "Warning", "Error", "unable to find SqlHost")
		return scheduledResult, client.IgnoreNotFound(err)
	}

	if getNamespacedName(database.Spec.HostRef, grants.Namespace) != getNamespacedName(user.Spec.HostRef, grants.Namespace) {
		err := fmt.Errorf("SqlDatabase and SqlUser don't share the same SqlHost")
		r.Recorder.Event(&grants, "Warning", "Error", err.Error())
		log.Error(err, "unable to find SqlHost for "+grants.Name)
		return scheduledResult, client.IgnoreNotFound(err)
	}

	driver, err := drivers.GetDriver(host)
	if err != nil {
		return ctrl.Result{}, err
	}

	finalizerName := "stenic.io/sqlgrants-deletion"
	// examine DeletionTimestamp to determine if object is under deletion
	if grants.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&grants, finalizerName) {
			controllerutil.AddFinalizer(&grants, finalizerName)
			if err := r.Update(ctx, &grants); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&grants, finalizerName) {
			// our finalizer is present, so lets handle any external dependency

			if grants.Spec.CleanupPolicy == steniciov1alpha1.CleanupPolicyDelete {
				// delete the user
				if err = driver.DeleteGrants(ctx, grants, user, database); err != nil {
					return ctrl.Result{}, err
				}
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&grants, finalizerName)
			if err := r.Update(ctx, &grants); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Deduplicate control loop
	if grants.Status.LastModifiedTimestamp != nil && time.Since(grants.Status.LastModifiedTimestamp.Time) < r.RefreshRate {
		return ctrl.Result{}, nil
	}

	if grants.Status.CurrentGrants == nil {
		grants.Status.CurrentGrants = []string{}
	}

	if !grants.Status.Created {
		grants.Status.Created = true
		grants.Status.CreationTimestamp = &metav1.Time{Time: time.Now()}
	}
	grants.Status.LastModifiedTimestamp = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &grants); err != nil {
		log.Error(err, "unable to update SqlGrant status")
		return ctrl.Result{}, err
	}

	count, err := driver.UpsertGrants(ctx, grants, user, database)
	if err != nil {
		log.Error(err, "failed to create SqlGrant")
		return ctrl.Result{}, err
	}
	if count > 0 {
		r.Recorder.Event(&grants, "Normal", "Changed", fmt.Sprintf("%d queries executed", count))
	}

	grants.Status.CurrentGrants = grants.Spec.Grants
	if err := r.Status().Update(ctx, &grants); err != nil {
		log.Error(err, "unable to update SqlGrant status")
		return ctrl.Result{}, err
	}

	return scheduledResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlGrantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlGrant{}).
		Complete(r)
}
