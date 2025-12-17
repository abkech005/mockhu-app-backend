#!/bin/bash
set -e

echo "🐳 Deploying Mockhu Backend to STAGING on GCP VM..."
echo "=================================================="

# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
echo "📦 Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo "✅ Docker installed"
else
    echo "✅ Docker already installed"
fi

# Install Docker Compose
echo "📦 Installing Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "✅ Docker Compose installed"
else
    echo "✅ Docker Compose already installed"
fi

# Verify installation
docker --version
docker-compose --version

echo ""
echo "=================================================="
echo "✅ STAGING Environment Setup Complete!"
echo "=================================================="
echo ""
echo "📝 Next steps to deploy:"
echo ""
echo "1. Clone your repository:"
echo "   git clone YOUR_REPO_URL /opt/mockhu-staging"
echo "   cd /opt/mockhu-staging"
echo ""
echo "2. Start STAGING environment:"
echo "   docker-compose -f docker-compose.staging.yml up -d"
echo ""
echo "3. View logs:"
echo "   docker-compose -f docker-compose.staging.yml logs -f"
echo ""
echo "4. Check status:"
echo "   docker-compose -f docker-compose.staging.yml ps"
echo ""
echo "5. Run migrations (if needed):"
echo "   docker-compose -f docker-compose.staging.yml exec backend ./mockhu-api migrate"
echo ""
echo "🌐 Your STAGING API will be available at:"
echo "   http://$(curl -s ifconfig.me)"
echo ""
echo "🔧 Useful commands:"
echo "   - Restart: docker-compose -f docker-compose.staging.yml restart"
echo "   - Stop: docker-compose -f docker-compose.staging.yml down"
echo "   - Rebuild: docker-compose -f docker-compose.staging.yml up -d --build"
echo ""
