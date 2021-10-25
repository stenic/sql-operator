#!/usr/bin/env bash
# Copyright 2021 Stenic BV.
# 
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# exit immediately when a command fails
set -e
# only exit with zero if all commands of the pipeline exit successfully
set -o pipefail
# error on unset variables
set -u

echo "Copying CRDs and Role to helm repo"

CONFIG_DIRECTORY="${CONFIG_DIRECTORY:-config}"
HELM_DIRECTORY="${HELM_DIRECTORY:-charts}"

# Copy and package CRDs
{
  cat hack/headers/header.yaml.txt
  printf "\n"
  cat ${CONFIG_DIRECTORY}/crd/bases/*.yaml
} > "${HELM_DIRECTORY}/sql-operator/crds/crds.yaml"

# Copy Kube Role for Operator permissions to Helm
# Template the Operator role as needed for Helm values
{
  cat hack/headers/header.yaml.txt
  printf '\n\n{{- if .Values.rbac.create }}\n{{- range $namespace := (split "," (include "sql-operator.watchNamespaces" $)) }}\n'
  cat ${CONFIG_DIRECTORY}/rbac/role.yaml \
    | awk '/^rules:$/{print "  namespace: {{ $namespace }}"}1' \
    | sed -E 's/^kind: ClusterRole$/kind: {{ include "sql-operator\.roleType" \$ }}/' \
    | sed -E 's/name: sql-operator-role$/name: {{ include "sql-operator\.fullname" \$ }}-role/'
  printf '\n{{- end }}\n\n---\n'
  cat ${CONFIG_DIRECTORY}/rbac/leader_election_role.yaml \
    | sed -E 's/name: leader-election-role$/name: {{ include "sql-operator\.fullname" \$ }}-leader-election-role/'
  printf '{{- end }}\n'
} > "${HELM_DIRECTORY}/sql-operator/templates/role.yaml"