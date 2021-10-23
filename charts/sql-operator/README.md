# sql-operator

## TL;DR;

```console
helm repo add stenic https://charts.stenic.io
helm install my-release --set "service.type=LoadBalancer" stenic/sql-operator
```

## Introduction

This chart installs `sql-operator` on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

[![openvpn-as](https://raw.githubusercontent.com/stenic/helm-charts/master/img/openvpn-as.png)](https://openvpn.net/index.php/access-server/overview.html)

## Prerequisites

- Kubernetes 1.12+
- Helm 3.0+
- LoadBalancer to expose the vpn service
- PV provisioner support in the underlying infrastructure

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add stenic https://charts.stenic.io
helm install my-release --set "service.type=LoadBalancer" stenic/sql-operator
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
| affinity | object | `{}` |  |
| annotations | object | `{}` |  |
| envVars | list | `[]` |  |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` |  |
| image.repository | string | `"ghcr.io/stenic/sql-operator"` |  |
| image.tag | string | `"1.0.0"` |  |
| labels | object | `{}` |  |
| nameOverride | string | `""` |  |
| nodeSelector | object | `{}` |  |
| priorityClassName | string | `""` |  |
| rbac.create | bool | `true` |  |
| replicaCount | int | `1` |  |
| resources | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `""` |  Required if create is false. If not set and create is true, a name is generated using the fullname template |
| sidecarContainers | list | `[]` |  |
| tolerations | list | `[]` |  |
| watchNamespaces | string | `""` |  If empty, the sql operator will watch all namespaces in the cluster. |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml stenic/sql-operator
```
