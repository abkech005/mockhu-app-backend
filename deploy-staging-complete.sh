#!/bin/bash
set -e

echo "🚀 Mockhu Backend - Complete Staging Deployment Script"
echo "======================================================="
echo ""

# Configuration
REPO_URL="https://github.com/abkech005/mockhu-app-backend.git"
DEPLOY_DIR="/opt/mockhu-staging"
COMPOSE_FILE="docker-compose.staging.yml"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}📦 Step 1: Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
    echo -e "${GREEN}✅ Docker installed${NC}"
else
    echo -e "${GREEN}✅ Docker already installed${NC}"
fi

echo ""
echo -e "${BLUE}📦 Step 2: Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
else
    echo -e "${GREEN}✅ Docker Compose already installed${NC}"
fi

# Verify installations
echo ""
echo -e "${BLUE}🔍 Verifying installations...${NC}"
docker --version
docker-compose --version

echo ""
echo -e "${BLUE}📥 Step 3: Cloning repository...${NC}"
if [ -d "$DEPLOY_DIR" ]; then
    echo -e "${YELLOW}⚠️  Directory $DEPLOY_DIR already exists. Updating...${NC}"
    cd $DEPLOY_DIR
    git pull origin main
else
    sudo mkdir -p $DEPLOY_DIR
    sudo chown $USER:$USER $DEPLOY_DIR
    git clone $REPO_URL $DEPLOY_DIR
    cd $DEPLOY_DIR
fi
echo -e "${GREEN}✅ Repository ready${NC}"

echo ""
echo -e "${BLUE}🐳 Step 4: Starting staging environment...${NC}"
# Stop existing containers if any
docker-compose -f $COMPOSE_FILE down 2>/dev/null || true

# Start all services
docker-compose -f $COMPOSE_FILE up -d

echo ""
echo -e "${BLUE}⏳ Waiting for services to start (10 seconds)...${NC}"
sleep 10

echo ""
echo -e "${BLUE}📊 Step 5: Checking container status...${NC}"
docker-compose -f $COMPOSE_FILE ps

echo ""
echo "======================================================="
echo -e "${GREEN}🎉 STAGING DEPLOYMENT COMPLETE!${NC}"
echo "======================================================="
echo ""
echo "📍 Deployment Location: $DEPLOY_DIR"
echo "🌐 API URL: http://$(curl -s ifconfig.me)"
echo ""
echo "🔧 Useful Commands:"
echo "  View logs:     docker-compose -f $COMPOSE_FILE logs -f"
echo "  Stop all:      docker-compose -f $COMPOSE_FILE down"
echo "  Restart:       docker-compose -f $COMPOSE_FILE restart"
echo "  Rebuild:       docker-compose -f $COMPOSE_FILE up -d --build"
echo ""
echo "📝 Next steps:"
echo "  1. Check logs: cd $DEPLOY_DIR && docker-compose -f $COMPOSE_FILE logs -f"
echo "  2. Test API: curl http://localhost/health"
echo "  3. Run migrations if needed"
echo ""
echo -e "${YELLOW}⚠️  Important: If you see permission errors, log out and back in, then re-run this script${NC}"
echo ""

# Show recent logs
echo -e "${BLUE}📋 Recent logs (last 20 lines):${NC}"
docker-compose -f $COMPOSE_FILE logs --tail=20

echo ""
echo "To view live logs, run:"
echo "  cd $DEPLOY_DIR && docker-compose -f $COMPOSE_FILE logs -f"
