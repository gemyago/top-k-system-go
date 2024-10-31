#!/usr/bin/env bash

set -e -o pipefail

# Print usage
usage() {
  echo "Usage: invoke-k8s-job.sh [-d] [-n <namespace>] [-p <pod-template-name>] -c <command>"
  echo "Description: Invoke command as a k8s job"
  echo "Options:"
  echo "  -c <command>            Required. Command to run in the container"
  echo "  -n <namespace>          Optional. Namespace to run the job in. Default: platform-services"
  echo "  -p <pod-template-name>  Optional. Pod template name to use. Default: top-k-system-go-cli-pod-template"
  echo "  -a <arg>                Optional. Arguments to pass to the command. Can be used multiple times"
  echo "  -d                      Optional. Dry run on the server"
}

# Reading command line arguments
while getopts ":n:p:c:a:d" opt; do
  case ${opt} in
    n )
      NAMESPACE=$OPTARG
      ;;
    p )
      POD_TEMPLATE_NAME=$OPTARG
      ;;
    c )
      COMMAND=$OPTARG
      ;;
    a )
      ARGS+=("$OPTARG")
      ;;
    d )
      echo "Dry run"
      DRY_RUN="--dry-run=server"
      ;;
    \? )
      usage
      exit 1
      ;;
  esac
done

# Convert ARGS array to a JSON array, make it empty array if it's empty
if [ ${#ARGS[@]} -eq 0 ]; then
  ARGS_JSON="[]"
else
  ARGS_JSON=$(printf '%s\n' "${ARGS[@]}" | jq -R . | jq -s .)
fi

# Check if the required -c option is provided
if [ -z "$COMMAND" ]; then
  echo "Error: -c <command> is required."
  usage
  exit 1
fi

NAMESPACE=${NAMESPACE:-platform-services}
POD_TEMPLATE_NAME=${POD_TEMPLATE_NAME:-top-k-system-go-cli-pod-template}

echo "Invoking command: ${COMMAND} in namespace: ${NAMESPACE} using pod template: ${POD_TEMPLATE_NAME}"

# Get the pod template
POD_TEMPLATE=$(kubectl get podtemplate ${POD_TEMPLATE_NAME} -n ${NAMESPACE} -o json | jq -r '.template')

# Updating the template, changing spec.restartPolicy to Never
# For some reason k8s will set it to Always when creating pod template
POD_TEMPLATE=$(echo ${POD_TEMPLATE} | jq '.spec.restartPolicy = "Never"')

# Add the command and args to the pod template
POD_TEMPLATE=$(echo ${POD_TEMPLATE} | jq ".spec.containers[0].command = [\"${COMMAND}\"]")
POD_TEMPLATE=$(echo ${POD_TEMPLATE} | jq --argjson args "${ARGS_JSON}" '.spec.containers[0].args = $args')

JOB_NAME="top-k-system-go-cli-job-${COMMAND}-$(date +%s)"

# If dry run is enabled, print the job template
if [ -n "${DRY_RUN}" ]; then
  echo "Dry run enabled. Job template:"
  echo ${POD_TEMPLATE} | jq .
fi

# Create a job based on the pod template
# We will keep the job for 5 minutes after it's finished
kubectl apply ${DRY_RUN} -n "${NAMESPACE}" -f - <<EOF
{
  "apiVersion": "batch/v1",
  "kind": "Job",
  "metadata": {
    "name": "${JOB_NAME}"
  },
  "spec": {
    "template": ${POD_TEMPLATE},
    "backoffLimit": 0,
    "ttlSecondsAfterFinished": 300
  }
}
EOF