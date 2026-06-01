terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project         = var.project_id
}

variable "project_id" {
  type        = string
  description = "The project ID to deploy the GCS bucket into"
}

variable "environment" {
  type        = string
  description = "Environment name (dev, staging, prod)"
}

resource "google_storage_bucket" "bucket" {
  name          = "${var.project_id}-application-bucket"
  location      = "US"
  force_destroy = true

  labels = {
    managed_by  = "foobar"
  }

  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }
}
