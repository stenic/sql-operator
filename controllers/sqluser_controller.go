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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stenic/sql-operator/drivers"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SqlUserReconciler reconciles a SqlUser object
type SqlUserReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Recorder    record.EventRecorder
	RefreshRate time.Duration
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers/finalizers,verbs=update
//+kubebuilder:rbac:groups=stenic.io,resources=sqlhosts,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *SqlUserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	promLabels := prometheus.Labels{
		"crd":       "sqlUser",
		"namespace": req.Namespace,
		"name":      req.Name,
	}
	sqlOperatorActions.With(promLabels).Inc()

	var user steniciov1alpha1.SqlUser
	if err := r.Get(ctx, req.NamespacedName, &user); err != nil {
		// log.Error(err, "unable to fetch SqlUser")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, getNamespacedName(user.Spec.HostRef, user.Namespace), &host); err != nil {
		log.Error(err, "unable to find SqlHost for "+user.Name)
		r.Recorder.Event(&user, "Warning", "Error", "unable to find SqlHost")
		sqlOperatorActionsFailures.With(promLabels).Inc()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	driver, err := drivers.GetDriver(host)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Provide ownership data
	err = driver.SetOwnershipData(ctx, drivers.OwnerShipData{
		Type:     drivers.OwnerShipTypeUser,
		Name:     user.Spec.Credentials.Username,
		Resource: req.NamespacedName.String(),
		OwnerID:  user.Status.OwnerID,
	})
	if err != nil {
		r.Recorder.Event(&host, "Warning", "Error", err.Error())
		return ctrl.Result{RequeueAfter: r.RefreshRate * 2}, err
	}

	scheduledResult := ctrl.Result{RequeueAfter: r.RefreshRate}

	finalizerName := "stenic.io/sqluser-deletion"
	// examine DeletionTimestamp to determine if object is under deletion
	if user.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&user, finalizerName) {
			controllerutil.AddFinalizer(&user, finalizerName)
			if err := r.Update(ctx, &user); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&user, finalizerName) {
			// our finalizer is present, so lets handle any external dependency

			var children steniciov1alpha1.SqlGrantList
			if err := isReferenced(ctx, r.Client, &children, referencedUserKey, &user); err != nil {
				r.Recorder.Event(&user, "Warning", "Error", err.Error())
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
				r.Recorder.Event(&user, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				// might have been faster than referenced object, reschedule.
				return scheduledResult, err
			}
			r.Recorder.Event(&user, "Normal", "Delete", "Validated no child objects")

			// Cleanup ref
			if err := driver.DeleteOwnerState(ctx); err != nil {
				log.Error(err, "unable to cleanup ownership")
				return ctrl.Result{}, err
			}
			r.Recorder.Event(&user, "Normal", "Delete", "Removed owner references")

			// delete the user
			if err = driver.DeleteUser(ctx, user); err != nil {
				r.Recorder.Event(&user, "Warning", "Error", err.Error())
				sqlOperatorActionsFailures.With(promLabels).Inc()
				return ctrl.Result{}, err
			}
			r.Recorder.Event(&user, "Normal", "Delete", "Deleted mysql object")

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&user, finalizerName)
			if err := r.Update(ctx, &user); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	if driver.Noop {
		r.Recorder.Event(&user, "Normal", "Noop", "Determined object is not owned")
		return ctrl.Result{}, nil
	}

	// Deduplicate control loop
	if user.Status.LastModifiedTimestamp != nil && time.Since(user.Status.LastModifiedTimestamp.Time) < r.RefreshRate {
		return ctrl.Result{}, nil
	}

	count, err := driver.UpsertUser(ctx, user)
	if err != nil {
		log.Error(err, "failed to create SqlUser")
		r.Recorder.Event(&user, "Warning", "Error", err.Error())
		sqlOperatorActionsFailures.With(promLabels).Inc()
		return ctrl.Result{}, err
	}
	if count > 0 {
		r.Recorder.Event(&user, "Normal", "Changed", fmt.Sprintf("%d queries executed", count))
		sqlOperatorQueries.With(promLabels).Add(float64(count))

		if !user.Status.Created {
			user.Status.Created = true
			user.Status.CreationTimestamp = &metav1.Time{Time: time.Now()}
		}
		user.Status.LastModifiedTimestamp = &metav1.Time{Time: time.Now()}
		user.Status.OwnerID = driver.GetOwnerID()

		if err := r.Status().Update(ctx, &user); err != nil {
			log.Error(err, "unable to update SqlUser status")
			return ctrl.Result{}, err
		}
	}

	return scheduledResult, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &steniciov1alpha1.SqlGrant{}, referencedUserKey, func(rawObj client.Object) []string {
		object := rawObj.(*steniciov1alpha1.SqlGrant)
		ns := object.Spec.UserRef.Namespace
		if ns == "" {
			ns = object.Namespace
		}

		return []string{ns + "/" + object.Spec.UserRef.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlUser{}).
		Complete(r)
}
