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
	"regexp"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var sqldatabaselog = logf.Log.WithName("sqldatabase-resource")

func (r *SqlDatabase) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-stenic-io-v1alpha1-sqldatabase,mutating=true,failurePolicy=fail,sideEffects=None,groups=stenic.io,resources=sqldatabases,verbs=create;update,versions=v1alpha1,name=msqldatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &SqlDatabase{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *SqlDatabase) Default() {
	sqldatabaselog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-stenic-io-v1alpha1-sqldatabase,mutating=false,failurePolicy=fail,sideEffects=None,groups=stenic.io,resources=sqldatabases,verbs=create;update,versions=v1alpha1,name=vsqldatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &SqlDatabase{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *SqlDatabase) ValidateCreate() error {
	sqldatabaselog.Info("validate create", "name", r.Name)

	return r.validateDatabase(nil)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *SqlDatabase) ValidateUpdate(old runtime.Object) error {
	sqldatabaselog.Info("validate update", "name", r.Name)

	return r.validateDatabase(old)
}

func (r *SqlDatabase) validateDatabase(old runtime.Object) error {
	var allErrs field.ErrorList
	if err := r.validateDatabaseName(); err != nil {
		allErrs = append(allErrs, err)
	}

	if old != nil {
		if err := r.validateDatabaseNameChanged(old); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "stenic.io", Kind: "SqlDatabase"},
		r.Name, allErrs)
}

func (r *SqlDatabase) validateDatabaseName() *field.Error {
	reg := regexp.MustCompile(`^[0-9,a-z,A-Z$_]+$`)
	if !reg.MatchString(r.Spec.DatabaseName) {
		return field.Invalid(
			field.NewPath("spec").Child("databaseName"),
			r.Spec.DatabaseName,
			"Invalid database name",
		)
	}

	return nil
}

func (r *SqlDatabase) validateDatabaseNameChanged(old runtime.Object) *field.Error {
	oldDb := old.(*SqlDatabase)

	if oldDb.Spec.DatabaseName != r.Spec.DatabaseName {
		return field.Invalid(
			field.NewPath("spec").Child("databaseName"),
			r.Spec.DatabaseName,
			"Field is immutable",
		)
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *SqlDatabase) ValidateDelete() error {
	sqldatabaselog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
