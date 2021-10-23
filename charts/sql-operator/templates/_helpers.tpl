{{/*
   * Copyright 2021 Stenic BV.
   * 
   * Licensed under the Apache License, Version 2.0 (the "License");
   * you may not use this file except in compliance with the License.
   * You may obtain a copy of the License at
   * 
   *     http://www.apache.org/licenses/LICENSE-2.0
   * 
   * Unless required by applicable law or agreed to in writing, software
   * distributed under the License is distributed on an "AS IS" BASIS,
   * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   * See the License for the specific language governing permissions and
   * limitations under the License.
   */}}

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "sql-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "sql-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "sql-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "sql-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "sql-operator.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ required "Must provide a serviceAccount.name if serviceAccount.create is set to false" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Get the namespaces to watch (empty if the operator should watch the entire cluster).
If .Values.watchNamespaces = true, then use the release namespace.
If .Values.watchNamespaces is a string, use it.
If .Values.watchNamespaces is empty or false, return empty.
*/}}
{{- define "sql-operator.watchNamespaces" -}}
{{- if .Values.watchNamespaces -}}
{{- if kindIs "bool" .Values.watchNamespaces -}}
{{ .Release.Namespace }}
{{- else -}}
{{ .Values.watchNamespaces }}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Determine whether to use ClusterRoles or Roles
*/}}
{{- define "sql-operator.roleType" -}}
{{- if .Values.watchNamespaces -}}
    Role
{{- else -}}
    ClusterRole
{{- end -}}
{{- end -}}
