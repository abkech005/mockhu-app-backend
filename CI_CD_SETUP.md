# GitHub Actions CI/CD Setup

This repository uses GitHub Actions to automatically deploy to GCP when code is pushed to the `develop` branch.

## Workflow

**File:** `.github/workflows/deploy-staging.yml`

**Trigger:** Push to `develop` branch

**What it does:**
1. ✅ Checks out code
2. ✅ Sets up Google Cloud SDK
3. ✅ SSHs into GCP VM (`mockhu-backend-dev`)
4. ✅ Pulls latest code from `develop` branch
5. ✅ Builds Go application
6. ✅ Restarts systemd service
7. ✅ Verifies deployment with health check

## Required Secrets

Add these secrets in **GitHub Settings > Secrets and variables > Actions**:

### 1. `GCP_PROJECT_ID`
Your GCP project ID:
```
mockhu-481519
```

### 2. `GCP_SA_KEY`
Service account key JSON for GitHub Actions to access GCP.

#### Creating the Service Account:

```bash
# 1. Create service account
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions Deployer"

# 2. Grant necessary permissions
gcloud projects add-iam-policy-binding mockhu-481519 \
  --member="serviceAccount:github-actions@mockhu-481519.iam.gserviceaccount.com" \
  --role="roles/compute.instanceAdmin.v1"

gcloud projects add-iam-policy-binding mockhu-481519 \
  --member="serviceAccount:github-actions@mockhu-481519.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# 3. Create and download key
gcloud iam service-accounts keys create github-actions-key.json \
  --iam-account=github-actions@mockhu-481519.iam.gserviceaccount.com

# 4. Copy the entire contents of github-actions-key.json
cat github-actions-key.json
```

Then paste the entire JSON content as the `GCP_SA_KEY` secret in GitHub.

## Testing the Workflow

1. Create `develop` branch if it doesn't exist:
   ```bash
   git checkout -b develop
   git push origin develop
   ```

2. Make any change and push to `develop`:
   ```bash
   git add .
   git commit -m "Test deployment"
   git push origin develop
   ```

3. Check the workflow status:
   - Go to **GitHub > Actions** tab
   - See the deployment progress

## Manual Deployment

If you need to deploy manually without GitHub Actions:

```bash
# SSH into VM
gcloud compute ssh mockhu-backend-dev --zone=us-central1-a

# Run updates
cd /opt/mockhu
git pull origin develop
/usr/local/go/bin/go build -o mockhu-api cmd/api/main.go
sudo systemctl restart mockhu
```

## Staging URL

**API:** http://104.198.232.235

**Health Check:** http://104.198.232.235/health

---

**Note:** The workflow will fail if the health check doesn't return "ok" status within 10 seconds of deployment.
