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
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"github.com/stenic/sql-operator/drivers"
)

// SqlDatabaseReconciler reconciles a SqlDatabase object
type SqlDatabaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	RefreshRate time.Duration
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases/finalizers,verbs=update
//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SqlDatabase object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SqlDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var database steniciov1alpha1.SqlDatabase
	if err := r.Get(ctx, req.NamespacedName, &database); err != nil {
		// log.Error(err, "unable to fetch SqlDatabase")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, getNamespacedName(database.Spec.HostRef, database.Namespace), &host); err != nil {
		log.Error(err, "unable to find SqlHost for "+database.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	driver, err := drivers.GetDriver(host)
	if err != nil {
		return ctrl.Result{}, err
	}

	finalizerName := "stenic.io/sqldatabase-deletion"
	// examine DeletionTimestamp to determine if object is under deletion
	if database.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&database, finalizerName) {
			controllerutil.AddFinalizer(&database, finalizerName)
			if err := r.Update(ctx, &database); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&database, finalizerName) {
			// our finalizer is present, so lets handle any external dependency

			if database.Spec.CleanupPolicy == steniciov1alpha1.CleanupPolicyDelete {
				// delete the user
				if err = driver.DeleteDatabase(ctx, database); err != nil {
					return ctrl.Result{}, err
				}
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&database, finalizerName)
			if err := r.Update(ctx, &database); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Deduplicate control loop
	if database.Status.LastModifiedTimestamp != nil && time.Since(database.Status.LastModifiedTimestamp.Time) < r.RefreshRate {
		return ctrl.Result{}, nil
	}

	scheduledResult := ctrl.Result{RequeueAfter: r.RefreshRate}

	if !database.Status.Created {
		database.Status.Created = true
		database.Status.CreationTimestamp = &metav1.Time{Time: time.Now()}
	}
	database.Status.LastModifiedTimestamp = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &database); err != nil {
		log.Error(err, "unable to update SqlDatabase status")
		return ctrl.Result{}, err
	}

	if err = driver.UpsertDatabase(ctx, database); err != nil {
		log.Error(err, "failed to create SqlDatabase")
		return ctrl.Result{}, err
	}

	return scheduledResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlDatabase{}).
		Complete(r)
}
