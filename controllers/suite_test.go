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
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/stenic/sql-operator/api/v1alpha1"
	steniciov1alpha1 "github.com/stenic/sql-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	err = steniciov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

}, 60)

var _ = Describe("Sql-operator", func() {

	const (
		SqlHostName     = "test-host"
		SqlDatabaseName = "test-host"
	)

	var validHost = &v1alpha1.SqlHost{
		TypeMeta: v1.TypeMeta{
			APIVersion: "stenic.io/v1alpha1",
			Kind:       "SqlHost",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      SqlHostName,
			Namespace: "default",
		},
		Spec: v1alpha1.SqlHostSpec{
			Engine: "Mysql",
			DSN:    "username:password@tcp(127.0.0.1:3306)/",
		},
	}

	var validDatabase = &v1alpha1.SqlDatabase{
		TypeMeta: v1.TypeMeta{
			APIVersion: "stenic.io/v1alpha1",
			Kind:       "SqlDatabase",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      SqlDatabaseName,
			Namespace: "default",
		},
		Spec: v1alpha1.SqlDatabaseSpec{
			HostRef: v1alpha1.SqlObjectRef{
				Name: SqlHostName,
			},
			DatabaseName: "test123",
		},
	}

	Context("SqlHost Validations", func() {
		ctx := context.Background()
		It("Should fail on lowercase mysql", func() {
			By("Creating a SqlHost", func() {
				var invalid = validHost.DeepCopy()
				invalid.Spec.Engine = "mysql"
				err := k8sClient.Create(ctx, invalid)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("SqlHost Crud", func() {
		ctx := context.Background()
		It("Should create without issue", func() {
			By("Creating a SqlHost", func() {
				Expect(k8sClient.Create(ctx, validHost.DeepCopy())).Should(Succeed())
			})
			By("Deleting a SqlHost", func() {
				Expect(k8sClient.Delete(ctx, validHost.DeepCopy())).Should(Succeed())
			})
		})
	})

	Context("SqlDB Crud", func() {
		ctx := context.Background()
		It("Should create without issue", func() {
			By("Creating a SqlHost", func() {
				Expect(k8sClient.Create(ctx, validDatabase.DeepCopy())).Should(Succeed())
			})
			By("Deleting a SqlHost", func() {
				Expect(k8sClient.Delete(ctx, validDatabase.DeepCopy())).Should(Succeed())
			})
		})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
