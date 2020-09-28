#!/bin/bash -e

PROJECT_ID=$(gcloud config get-value project)
PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format="value(projectNumber)")
BUILD_SA=$PROJECT_NUMBER@cloudbuild.gserviceaccount.com
COMPUTE_SA=$PROJECT_NUMBER-compute@developer.gserviceaccount.com

# Allow Cloud Build to deploy Cloud Run services
gcloud projects add-iam-policy-binding \
  $PROJECT_ID \
  --member serviceAccount:$BUILD_SA \
  --role roles/run.admin

# Allow Cloud Build to use the Compute Engine default service account.
gcloud iam service-accounts add-iam-policy-binding $COMPUTE_SA \
  --member serviceAccount:$BUILD_SA \
  --role roles/iam.serviceAccountUser

# Allow yourself to impersonate the Compute Engine default service account.
gcloud iam service-accounts add-iam-policy-binding $COMPUTE_SA \
  --member user:$(gcloud config list account --format "value(core.account)") \
  --role roles/iam.serviceAccountTokenCreator