#!/bin/sh
set -e
echo "Environment variables:"
env

echo "Initializing Terraform..."
cd terraform/
terraform init -input=false -backend-config="bucket=app1-cicd-8c0e2f7a-tfstate" -backend-config="prefix=terraform/state/${environment}"

echo "Running Terraform Apply..."
terraform apply -lock=false -input=false -auto-approve -var="project_id=${project_id}" -var="environment=${environment}"

echo "{\"resultStatus\": \"SUCCEEDED\"}" > results.json
gsutil cp results.json "$CLOUD_DEPLOY_OUTPUT_GCS_PATH/results.json"
