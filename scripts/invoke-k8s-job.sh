#!/usr/bin/env bash

set -e -o pipefail

# This script is used to invoke k8s jobs based on the top-k-system-go-cli-pod-template PodTemplate 
# It allows invoking any command in the container

NAMESPACE=platform-services
POD_TEMPLATE_NAME=top-k-system-go-cli-pod-template

# Get the pod template
POD_TEMPLATE=$(kubectl get podtemplate ${POD_TEMPLATE_NAME} -n ${NAMESPACE} -o json | jq -r '.template')

# Updating the template, changing spec.restartPolicy to Never
# For some reason k8s will set it to Always when creating pod template
POD_TEMPLATE=$(echo ${POD_TEMPLATE} | jq '.spec.restartPolicy = "Never"')

JOB_NAME="top-k-system-go-cli-job-$(date +%s)"

# Create a job based on the pod template
kubectl apply -n "${NAMESPACE}" -f - <<EOF
{
  "apiVersion": "batch/v1",
  "kind": "Job",
  "metadata": {
    "name": "${JOB_NAME}"
  },
  "spec": {
    "template": ${POD_TEMPLATE},
    "backoffLimit": 0
  }
}
EOF