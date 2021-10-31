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

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var sqluserlog = logf.Log.WithName("sqluser-resource")

func (r *SqlUser) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-stenic-io-v1alpha1-sqluser,mutating=true,failurePolicy=fail,sideEffects=None,groups=stenic.io,resources=sqlusers,verbs=create;update,versions=v1alpha1,name=msqluser.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SqlUser{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SqlUser) Default() {
	sqluserlog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-stenic-io-v1alpha1-sqluser,mutating=false,failurePolicy=fail,sideEffects=None,groups=stenic.io,resources=sqlusers,verbs=create;update,versions=v1alpha1,name=vsqluser.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SqlUser{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SqlUser) ValidateCreate() error {
	sqluserlog.Info("validate create", "name", r.Name)

	return r.validateUser(nil)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SqlUser) ValidateUpdate(old runtime.Object) error {
	sqluserlog.Info("validate update", "name", r.Name)

	return r.validateUser(old)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SqlUser) ValidateDelete() error {
	sqluserlog.Info("validate delete", "name", r.Name)

	return nil
}

// validateUser runs all checks to validate the object
func (r *SqlUser) validateUser(old runtime.Object) error {
	var allErrs field.ErrorList

	if old != nil {
		if err := r.validateUsernameChanged(old); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "stenic.io", Kind: "SqlUser"},
		r.Name, allErrs)
}

// validateUsernameChanged ensures the username can't be changed
func (r *SqlUser) validateUsernameChanged(old runtime.Object) *field.Error {
	oldDb := old.(*SqlUser)

	if oldDb.Spec.Credentials.Username != r.Spec.Credentials.Username {
		return field.Invalid(
			field.NewPath("spec").Child("credentials").Child("username"),
			r.Spec.Credentials.Username,
			"Field is immutable",
		)
	}

	return nil
}
