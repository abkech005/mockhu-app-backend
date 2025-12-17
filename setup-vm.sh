#!/bin/bash
set -e

echo "🚀 Setting up Mockhu Backend on GCP VM..."
echo "=========================================="

# Update system
echo "📦 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install Go 1.21+
echo "🔧 Installing Go 1.21.5..."
cd /tmp
wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin
go version

# Install PostgreSQL 15
echo "🗄️ Installing PostgreSQL 15..."
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update
sudo apt install -y postgresql-15 postgresql-contrib-15

sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create database and user
echo "🔐 Creating database and user..."
sudo -u postgres psql << EOF
CREATE DATABASE mockhu_db;
CREATE USER mockhu WITH ENCRYPTED PASSWORD 'mockhu_dev_password_2024';
GRANT ALL PRIVILEGES ON DATABASE mockhu_db TO mockhu;
ALTER DATABASE mockhu_db OWNER TO mockhu;
\q
EOF

echo "✅ PostgreSQL setup complete"

# Install Nginx
echo "🌐 Installing Nginx..."
sudo apt install -y nginx git curl wget make
sudo systemctl enable nginx
sudo systemctl start nginx

# Clone Mockhu repository
echo "📂 Cloning Mockhu repository..."
sudo mkdir -p /opt/mockhu
sudo chown $USER:$USER /opt/mockhu

# You'll need to add your authentication here or use a personal access token
read -p "Enter your GitHub username: " GITHUB_USER
echo "Clone your repository manually or continue with public clone..."
# Uncomment and modify: git clone https://github.com/$GITHUB_USER/mockhu-app-backend.git /opt/mockhu

echo "⚠️  Please clone your repository manually:"
echo "   cd /opt/mockhu"
echo "   git clone https://github.com/YOUR_USERNAME/mockhu-app-backend.git ."

# Create .env file
echo "⚙️ Creating .env file..."
sudo mkdir -p /opt/mockhu
cat > /tmp/.env << 'ENVEOF'
# Database Configuration
DATABASE_URL=postgresql://mockhu:mockhu_dev_password_2024@localhost:5432/mockhu_db

# JWT Configuration
JWT_SECRET=dev_jwt_secret_change_in_production_2024

# Server Configuration
PORT=8085
ENV=development

# CORS (if needed)
ALLOWED_ORIGINS=*
ENVEOF

sudo mv /tmp/.env /opt/mockhu/.env

# Create systemd service
echo "⚙️ Creating systemd service..."
sudo tee /etc/systemd/system/mockhu.service > /dev/null << 'SERVICEEOF'
[Unit]
Description=Mockhu Backend API
After=network.target postgresql.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/mockhu
ExecStart=/opt/mockhu/mockhu-api
Restart=always
RestartSec=5s
Environment="PATH=/usr/local/go/bin:/usr/bin:/bin"
StandardOutput=journal
StandardError=journal
SyslogIdentifier=mockhu

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Configure Nginx
echo "🌐 Configuring Nginx..."
sudo tee /etc/nginx/sites-available/mockhu > /dev/null << 'NGINXEOF'
server {
    listen 80;
    server_name _;
    client_max_body_size 100M;

    location / {
        proxy_pass http://localhost:8085;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    location /health {
        proxy_pass http://localhost:8085/health;
        access_log off;
    }

    location /avatars/ {
        alias /opt/mockhu/storage/avatars/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
NGINXEOF

sudo ln -sf /etc/nginx/sites-available/mockhu /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx

# Create storage directories
echo "📁 Creating storage directories..."
sudo mkdir -p /opt/mockhu/storage/avatars
sudo mkdir -p /opt/mockhu/storage/posts
sudo mkdir -p /opt/mockhu/storage/messages
sudo chown -R $USER:$USER /opt/mockhu/storage

echo ""
echo "=========================================="
echo "✅ Initial setup complete!"
echo "=========================================="
echo ""
echo "📝 Next steps:"
echo "1. Clone your repository:"
echo "   cd /opt/mockhu"
echo "   git clone YOUR_REPO_URL ."
echo ""
echo "2. Run database migrations:"
echo "   cd /opt/mockhu"
echo "   for migration in migrations/*.up.sql; do"
echo "     sudo -u postgres psql -d mockhu_db -f \"\$migration\""
echo "   done"
echo ""
echo "3. Build the application:"
echo "   cd /opt/mockhu"
echo "   go mod download"
echo "   go build -o mockhu-api cmd/api/main.go"
echo ""
echo "4. Start the service:"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl enable mockhu"
echo "   sudo systemctl start mockhu"
echo ""
echo "5. Check status:"
echo "   sudo systemctl status mockhu"
echo "   sudo journalctl -u mockhu -f"
echo ""
echo "🌐 Your API will be available at: http://34.30.185.153"
echo ""
