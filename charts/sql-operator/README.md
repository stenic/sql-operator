# sql-operator

## TL;DR;

```console
helm repo add sql-operator https://stenic.github.io/sql-operator/
helm install sql-operator --namespace sql-operator sql-operator/sql-operator
```

## Introduction

This chart installs `sql-operator` on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.12+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add sql-operator https://stenic.github.io/sql-operator/
helm install sql-operator --namespace sql-operator sql-operator/sql-operator
```

These commands deploy sql-operator on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables list the configurable parameters of the sql-operator chart and their default values.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity and anti-affinity |
| annotations | object | `{}` | Additional annotations for the controller pods. |
| controller.logLevel | string | `"info"` | Set the to configure the verbosity of logging. Can be one of 'debug', 'info', 'error'. |
| envVars | list | `[]` | Additional environment variables for the controller. |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/stenic/sql-operator"` |  |
| image.tag | string | `""` |  |
| labels | object | `{}` | Additional labels for the controller pods. |
| nodeSelector | object | `{"kubernetes.io/os":"linux"}` | Node labels for controller pod assignment |
| priorityClassName | string | `""` | Provide a priority class name to the controller pods |
| rbac.create | bool | `true` | Specifies whether RBAC resources should be created |
| replicaCount | int | `1` |  |
| resources | object | `{}` | Resource requests and limits for the controller |
| serviceAccount.create | bool | `true` | Specifies whether a ServiceAccount should be created |
| serviceAccount.name | string | `""` | The name of the ServiceAccount to use. Required if create is false. If not set and create is true, a name is generated using the fullname template |
| sidecarContainers | list | `[]` | Additional containers to be added to the controller pod. |
| tolerations | list | `[]` | Node tolerations for server scheduling to nodes with taints |
| watchNamespaces | string | `""` | A comma-separated list of namespaces that the operator should watch. If empty, the sql operator will watch all namespaces in the cluster. |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml sql-operator/sql-operator
```
