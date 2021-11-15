package controllers

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	sqlOperatorActions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sqloperator_actions_total",
			Help: "Number of sql-operator actions proccessed",
		},
		[]string{"crd", "namespace", "name"},
	)
	sqlOperatorActionsFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sqloperator_actions_failures_total",
			Help: "Number of sql-operator actions failures",
		},
		[]string{"crd", "namespace", "name"},
	)
	sqlOperatorQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sqloperator_queries_total",
			Help: "Number of sql-operator queries executed",
		},
		[]string{"crd", "namespace", "name"},
	)
)

func init() {
	metrics.Registry.MustRegister(
		sqlOperatorActions,
		sqlOperatorActionsFailures,
		sqlOperatorQueries,
	)
}
