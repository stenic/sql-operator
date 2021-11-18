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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/prometheus/client_golang/prometheus"
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	"github.com/stenic/sql-operator/drivers"
)

// SqlHostReconciler reconciles a SqlHost object
type SqlHostReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Recorder    record.EventRecorder
	RefreshRate time.Duration
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SqlHostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	promLabels := prometheus.Labels{
		"crd":       "sqlHost",
		"namespace": req.Namespace,
		"name":      req.Name,
	}
	sqlOperatorActions.With(promLabels).Inc()

	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, req.NamespacedName, &host); err != nil {
		// log.Error(err, "unable to fetch SqlDatabase")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	driver, err := drivers.GetDriver(host)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := driver.InitOwnerSchema(ctx); err != nil {
		r.Recorder.Event(&host, "Warning", "Error", err.Error())
		return ctrl.Result{RequeueAfter: r.RefreshRate * 10}, err
	}

	scheduledResult := ctrl.Result{RequeueAfter: r.RefreshRate}

	finalizerName := "stenic.io/sqlhost-deletion"
	// examine DeletionTimestamp to determine if object is under deletion
	if host.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&host, finalizerName) {
			controllerutil.AddFinalizer(&host, finalizerName)
			if err := r.Update(ctx, &host); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&host, finalizerName) {
			// our finalizer is present, so lets handle any external dependency

			var userChildren steniciov1alpha1.SqlUserList
			if err := isReferenced(ctx, r.Client, &userChildren, referencedHostKey, &host); err != nil {
				r.Recorder.Event(&host, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				return ctrl.Result{}, err
			}
			if len(userChildren.Items) > 0 {
				err := fmt.Errorf(
					"%s - [%s/%s] ...",
					"Can't delete, found other referencing this object",
					userChildren.Items[0].Namespace,
					userChildren.Items[0].Name,
				)
				r.Recorder.Event(&host, "Warning", "Error", err.Error())
				// might have been faster than referenced object, reschedule.
				return scheduledResult, err
			}

			var databaseChildren steniciov1alpha1.SqlDatabaseList
			if err := isReferenced(ctx, r.Client, &databaseChildren, referencedHostKey, &host); err != nil {
				r.Recorder.Event(&host, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				return ctrl.Result{}, err
			}
			if len(databaseChildren.Items) > 0 {
				err := fmt.Errorf(
					"%s - [%s/%s] ...",
					"Can't delete, found other referencing this object",
					databaseChildren.Items[0].Namespace,
					databaseChildren.Items[0].Name,
				)
				r.Recorder.Event(&host, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				// might have been faster than referenced object, reschedule.
				return scheduledResult, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&host, finalizerName)
			if err := r.Update(ctx, &host); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlHostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &steniciov1alpha1.SqlUser{}, referencedHostKey, func(rawObj client.Object) []string {
		object := rawObj.(*steniciov1alpha1.SqlUser)
		ns := object.Spec.HostRef.Namespace
		if ns == "" {
			ns = object.Namespace
		}

		return []string{ns + "/" + object.Spec.HostRef.Name}
	}); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &steniciov1alpha1.SqlDatabase{}, referencedHostKey, func(rawObj client.Object) []string {
		object := rawObj.(*steniciov1alpha1.SqlDatabase)
		ns := object.Spec.HostRef.Namespace
		if ns == "" {
			ns = object.Namespace
		}

		return []string{ns + "/" + object.Spec.HostRef.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlHost{}).
		Complete(r)
}
