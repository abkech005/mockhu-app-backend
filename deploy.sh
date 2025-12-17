#!/bin/bash

# Mockhu Staging Deployment Script
# Run this from your local machine to deploy to GCP VM

set -e

echo "🚀 Deploying to Mockhu Staging..."
echo "=================================="
echo ""

# Configuration
VM_NAME="mockhu-backend-dev"
ZONE="us-central1-a"
BRANCH="${1:-main}"  # Default to main branch, or use argument

echo "📍 VM: $VM_NAME"
echo "🌍 Zone: $ZONE"
echo "🌿 Branch: $BRANCH"
echo ""

# Deploy via SSH
echo "🔄 Connecting to VM and deploying..."
gcloud compute ssh $VM_NAME --zone=$ZONE --command="
  set -e
  echo '📥 Pulling latest code from $BRANCH...'
  cd /opt/mockhu
  git pull origin $BRANCH
  
  echo '🔨 Building application...'
  /usr/local/go/bin/go build -o mockhu-api cmd/api/main.go
  
  echo '🔄 Restarting service...'
  sudo systemctl restart mockhu
  
  echo '⏳ Waiting for service to start...'
  sleep 3
  
  echo '✅ Deployment complete!'
  echo ''
  echo '📊 Service Status:'
  sudo systemctl status mockhu --no-pager -l | head -20
"

echo ""
echo "=================================="
echo "🎉 Deployment Successful!"
echo "=================================="
echo ""
echo "🌐 Staging URL: http://104.198.232.235"
echo "🏥 Health Check: http://104.198.232.235/health"
echo ""
echo "📝 Quick commands:"
echo "  View logs:    gcloud compute ssh $VM_NAME --zone=$ZONE --command='sudo journalctl -u mockhu -f'"
echo "  Check status: gcloud compute ssh $VM_NAME --zone=$ZONE --command='sudo systemctl status mockhu'"
echo ""
