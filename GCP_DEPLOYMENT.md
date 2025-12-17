# 🚀 Mockhu Backend - GCP Deployment

## Quick Deploy to GCP (Native - No Docker)

### Prerequisites
- GCP account with billing enabled
- `gcloud` CLI installed and authenticated

### One-Command Deployment

```bash
# 1. Create VM
gcloud compute instances create mockhu-backend-dev \
  --project=YOUR_PROJECT_ID \
  --zone=us-central1-a \
  --machine-type=e2-small \
  --image-family=ubuntu-2204-lts \
  --image-project=ubuntu-os-cloud \
  --boot-disk-size=20GB \
  --boot-disk-type=pd-ssd \
  --tags=http-server,https-server

# 2. SSH into VM
gcloud compute ssh mockhu-backend-dev --zone=us-central1-a

# 3. Run deployment script (on VM)
curl -fsSL https://raw.githubusercontent.com/abkech005/mockhu-app-backend/main/deploy-native.sh | bash
```

### What Gets Installed
- PostgreSQL 15
- Go 1.24
- Nginx (reverse proxy)
- Your backend app
- Systemd service (auto-restart)

### Cost Estimate
- **Development:** ~$15-20/month (e2-small + 20GB SSD)
- **Production:** ~$200-350/month (with load balancing, managed DB)

### Useful Commands

```bash
# Check service status
sudo systemctl status mockhu

# View logs
sudo journalctl -u mockhu -f

# Restart service
sudo systemctl restart mockhu

# Update code
cd /opt/mockhu
git pull origin main
go build -o mockhu-api cmd/api/main.go
sudo systemctl restart mockhu
```

### Staging URLs
- **Current Staging:** http://104.198.232.235
- **VM Name:** mockhu-backend-dev
- **Zone:** us-central1-a

---

For detailed deployment guides, see:
- [Full GCP Guide](/Users/abkech/.gemini/antigravity/brain/0da0048c-8180-49e0-ba62-bc1bc220fff0/gcp_deployment_guide.md)
- [Development Setup](/Users/abkech/.gemini/antigravity/brain/0da0048c-8180-49e0-ba62-bc1bc220fff0/gcp_dev_deployment.md)
