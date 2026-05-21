#!/bin/sh
set -e
echo "Environment variables:"
env

echo "Initializing Terraform..."
cd terraform/
terraform init -input=false -backend-config="bucket=cicd-project-a9403b81-tfstate" -backend-config="prefix=terraform/state/${environment}"

echo "Running Terraform Plan..."
terraform plan -lock=false -input=false -var="project_id=${project_id}" -var="environment=${environment}"

echo '' > manifest.yaml
gsutil cp manifest.yaml "$CLOUD_DEPLOY_OUTPUT_GCS_PATH/manifest.yaml"

echo "{\"resultStatus\": \"SUCCEEDED\", \"manifestFile\": \"$CLOUD_DEPLOY_OUTPUT_GCS_PATH/manifest.yaml\"}" > results.json
gsutil cp results.json "$CLOUD_DEPLOY_OUTPUT_GCS_PATH/results.json"

cd /workspace
tar -czvf workspace.tgz .
gsutil cp workspace.tgz "$CLOUD_DEPLOY_OUTPUT_GCS_PATH/workspace.tgz"
