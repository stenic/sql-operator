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

	"github.com/prometheus/client_golang/prometheus"
)

// SqlDatabaseReconciler reconciles a SqlDatabase object
type SqlDatabaseReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	RefreshRate time.Duration
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqldatabases/finalizers,verbs=update
//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SqlDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	promLabels := prometheus.Labels{
		"crd":       "sqlDatabase",
		"namespace": req.Namespace,
		"name":      req.Name,
	}

	sqlOperatorActions.With(promLabels).Inc()

	var database steniciov1alpha1.SqlDatabase
	if err := r.Get(ctx, req.NamespacedName, &database); err != nil {
		// log.Error(err, "unable to fetch SqlDatabase")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, getNamespacedName(database.Spec.HostRef, database.Namespace), &host); err != nil {
		log.Error(err, "unable to find SqlHost for "+database.Name)
		r.Recorder.Event(&database, "Warning", "Error", "unable to find SqlHost")
		sqlOperatorActionsFailures.With(promLabels).Inc()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	driver, err := drivers.GetDriver(host)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Provide ownership data
	err = driver.SetOwnershipData(ctx, drivers.OwnerShipData{
		Type:     drivers.OwnerShipTypeDatabase,
		Name:     database.Spec.DatabaseName,
		Resource: req.NamespacedName.String(),
		OwnerID:  database.Status.OwnerID,
	})
	if err != nil {
		r.Recorder.Event(&host, "Warning", "Error", err.Error())
		return ctrl.Result{RequeueAfter: r.RefreshRate * 2}, err
	}

	scheduledResult := ctrl.Result{RequeueAfter: r.RefreshRate}

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

			var children steniciov1alpha1.SqlGrantList
			if err := isReferenced(ctx, r.Client, &children, referencedDatabaseKey, &database); err != nil {
				r.Recorder.Event(&database, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				return ctrl.Result{}, err
			}
			if len(children.Items) > 0 {
				err := fmt.Errorf(
					"%s - [%s/%s] ...",
					"can't delete, found other referencing this object",
					children.Items[0].Namespace,
					children.Items[0].Name,
				)
				r.Recorder.Event(&database, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				// might have been faster than referenced object, reschedule.
				return scheduledResult, err
			}
			r.Recorder.Event(&database, "Normal", "Delete", "Validated no child objects")

			// Cleanup ref
			if err := driver.DeleteOwnerState(ctx); err != nil {
				log.Error(err, "unable to cleanup ownership")
				return ctrl.Result{}, err
			}
			r.Recorder.Event(&database, "Normal", "Delete", "Removed owner references")

			// delete the user
			if err = driver.DeleteDatabase(ctx, database); err != nil {
				r.Recorder.Event(&database, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				return ctrl.Result{}, err
			}
			r.Recorder.Event(&database, "Normal", "Delete", "Deleted mysql object")

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&database, finalizerName)
			if err := r.Update(ctx, &database); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	if driver.Noop {
		r.Recorder.Event(&database, "Normal", "Noop", "Determined object is not owned")
		return ctrl.Result{}, nil
	}

	// Deduplicate control loop
	if database.Status.LastModifiedTimestamp != nil && time.Since(database.Status.LastModifiedTimestamp.Time) < r.RefreshRate {
		return ctrl.Result{}, nil
	}

	count, err := driver.UpsertDatabase(ctx, database)
	if err != nil {
		log.Error(err, "failed to create SqlDatabase")
		r.Recorder.Event(&database, "Warning", "Error", err.Error())
		sqlOperatorActionsFailures.With(promLabels).Inc()
		return ctrl.Result{}, err
	}
	if count > 0 {
		r.Recorder.Event(&database, "Normal", "Changed", fmt.Sprintf("%d queries executed", count))
		sqlOperatorQueries.With(promLabels).Add(float64(count))

		if !database.Status.Created {
			database.Status.Created = true
			database.Status.CreationTimestamp = &metav1.Time{Time: time.Now()}
		}
		database.Status.LastModifiedTimestamp = &metav1.Time{Time: time.Now()}
		database.Status.OwnerID = driver.GetOwnerID()

		if err := r.Status().Update(ctx, &database); err != nil {
			log.Error(err, "unable to update SqlDatabase status")
			return ctrl.Result{}, err
		}
	}

	return scheduledResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &steniciov1alpha1.SqlGrant{}, referencedDatabaseKey, func(rawObj client.Object) []string {
		object := rawObj.(*steniciov1alpha1.SqlGrant)
		ns := object.Spec.DatabaseRef.Namespace
		if ns == "" {
			ns = object.Namespace
		}

		return []string{ns + "/" + object.Spec.DatabaseRef.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlDatabase{}).
		Complete(r)
}
