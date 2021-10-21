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
	"k8s.io/apimachinery/pkg/types"
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
}

//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stenic.io,resources=sqlusers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SqlUser object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SqlUserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var user steniciov1alpha1.SqlUser
	if err := r.Get(ctx, req.NamespacedName, &user); err != nil {
		// log.Error(err, "unable to fetch SqlUser")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	finalizerName := "stenic.io/sqluser-deletion"

	// examine DeletionTimestamp to determine if object is under deletion
	if user.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(user.GetFinalizers(), finalizerName) {
			controllerutil.AddFinalizer(&user, finalizerName)
			if err := r.Update(ctx, &user); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(user.GetFinalizers(), finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			hostNamespacedName := r.getHostNamespacedName(user)
			var host steniciov1alpha1.SqlHost
			if err := r.Get(ctx, hostNamespacedName, &host); err != nil {
				log.Error(err, "unable to find SqlHost "+hostNamespacedName.Name+" in "+hostNamespacedName.Namespace)
				return ctrl.Result{}, err
			}

			// delete the user
			if err := r.deleteExternalResource(ctx, user, host); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&user, finalizerName)
			if err := r.Update(ctx, &user); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	hostNamespacedName := r.getHostNamespacedName(user)
	var host steniciov1alpha1.SqlHost
	if err := r.Get(ctx, hostNamespacedName, &host); err != nil {
		log.Error(err, "unable to find SqlHost "+hostNamespacedName.Name+" in "+hostNamespacedName.Namespace)
		return ctrl.Result{}, err
	}

	if !user.Status.Created {
		if err := r.createExternalResource(ctx, user, host); err != nil {
			log.Error(err, "failed to create SqlUser")
			return ctrl.Result{}, err
		}

		user.Status.Created = true
		user.Status.CreationTimestamp = &metav1.Time{Time: time.Now()}
	} else {
		if err := r.updateExternalResource(ctx, user, host); err != nil {
			log.Error(err, "failed to update SqlUser")
			return ctrl.Result{}, err
		}
	}

	user.Status.LastModifiedTimestamp = &metav1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &user); err != nil {
		log.Error(err, "unable to update SqlUser status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

func (r *SqlUserReconciler) createExternalResource(ctx context.Context, user steniciov1alpha1.SqlUser, host steniciov1alpha1.SqlHost) error {
	log := log.FromContext(ctx)

	log.Info(fmt.Sprintf(
		"CREATE USER %s WITH PASSWORD %s;",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	))

	return nil
}
func (r *SqlUserReconciler) deleteExternalResource(ctx context.Context, user steniciov1alpha1.SqlUser, host steniciov1alpha1.SqlHost) error {
	log := log.FromContext(ctx)

	log.Info(fmt.Sprintf(
		"DELETE USER %s WITH PASSWORD %s;",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	))

	return nil
}
func (r *SqlUserReconciler) updateExternalResource(ctx context.Context, user steniciov1alpha1.SqlUser, host steniciov1alpha1.SqlHost) error {
	log := log.FromContext(ctx)

	log.Info(fmt.Sprintf(
		"UPDATE USER %s WITH PASSWORD %s;",
		user.Spec.Credentials.Username,
		user.Spec.Credentials.Password,
	))

	return nil
}

func (r *SqlUserReconciler) getHostNamespacedName(user steniciov1alpha1.SqlUser) types.NamespacedName {
	if user.Spec.HostRef.Namespace != "" {
		return types.NamespacedName{
			Namespace: user.Spec.HostRef.Namespace,
			Name:      user.Spec.HostRef.Name,
		}
	}
	return types.NamespacedName{
		Namespace: user.Namespace,
		Name:      user.Spec.HostRef.Name,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SqlUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&steniciov1alpha1.SqlUser{}).
		Complete(r)
}
