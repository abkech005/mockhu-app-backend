#!/bin/bash
set -e

echo "🚀 Mockhu Staging - Native Installation (No Docker)"
echo "====================================================="
echo ""

# Configuration
REPO_URL="https://github.com/abkech005/mockhu-app-backend.git"
APP_DIR="/opt/mockhu"
DB_NAME="mockhu_staging_db"
DB_USER="mockhu"
DB_PASS="mockhu_staging_2024"

echo "📦 Step 1: Updating system packages..."
sudo apt update && sudo apt upgrade -y

echo ""
echo "🔧 Step 2: Installing Go 1.24..."
cd /tmp
wget -q https://go.dev/dl/go1.24rc2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24rc2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
export PATH=$PATH:/usr/local/go/bin
go version

echo ""
echo "🗄️ Step 3: Installing PostgreSQL 15..."
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update
sudo apt install -y postgresql-15 postgresql-contrib-15

sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create database and user
echo "🔐 Creating database..."
sudo -u postgres psql << EOF
DROP DATABASE IF EXISTS $DB_NAME;
DROP USER IF EXISTS $DB_USER;
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASS';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
ALTER DATABASE $DB_NAME OWNER TO $DB_USER;
\q
EOF

echo "✅ Database created: $DB_NAME"

echo ""
echo "🌐 Step 4: Installing Nginx..."
sudo apt install -y nginx git curl wget make
sudo systemctl enable nginx
sudo systemctl start nginx

echo ""
echo "📥 Step 5: Cloning repository..."
if [ -d "$APP_DIR" ]; then
    echo "Directory exists, updating..."
    cd $APP_DIR
    sudo git pull origin main
else
    sudo git clone $REPO_URL $APP_DIR
fi

cd $APP_DIR
sudo chown -R $USER:$USER $APP_DIR

echo ""
echo "⚙️ Step 6: Creating .env file..."
cat > $APP_DIR/.env << ENVEOF
DATABASE_URL=postgresql://$DB_USER:$DB_PASS@localhost:5432/$DB_NAME
JWT_SECRET=staging_jwt_secret_2024
PORT=8085
ENV=staging
ENVEOF

echo ""
echo "🔄 Step 7: Running database migrations..."
for migration in $APP_DIR/migrations/*.up.sql; do
    if [ -f "$migration" ]; then
        echo "Running: $migration"
        sudo -u postgres psql -d $DB_NAME -f "$migration" || true
    fi
done

echo ""
echo "🔨 Step 8: Building Go application..."
cd $APP_DIR
/usr/local/go/bin/go mod download
/usr/local/go/bin/go build -o mockhu-api cmd/api/main.go

echo ""
echo "📁 Step 9: Creating storage directories..."
mkdir -p $APP_DIR/storage/avatars
mkdir -p $APP_DIR/storage/posts
mkdir -p $APP_DIR/storage/messages

echo ""
echo "⚙️ Step 10: Creating systemd service..."
sudo tee /etc/systemd/system/mockhu.service > /dev/null << 'SERVICEEOF'
[Unit]
Description=Mockhu Backend API (Staging)
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

sudo systemctl daemon-reload
sudo systemctl enable mockhu
sudo systemctl start mockhu

echo ""
echo "🌐 Step 11: Configuring Nginx..."
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
    }

    location /storage/ {
        alias /opt/mockhu/storage/;
        expires 30d;
    }
}
NGINXEOF

sudo ln -sf /etc/nginx/sites-available/mockhu /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx

echo ""
echo "====================================================="
echo "🎉 STAGING DEPLOYMENT COMPLETE!"
echo "====================================================="
echo ""
echo "📊 Service Status:"
sudo systemctl status mockhu --no-pager -l
echo ""
echo "🌐 Your API is available at:"
echo "   http://$(curl -s ifconfig.me)"
echo ""
echo "🔧 Useful Commands:"
echo "   - View logs:        sudo journalctl -u mockhu -f"
echo "   - Restart service:  sudo systemctl restart mockhu"
echo "   - Check status:     sudo systemctl status mockhu"
echo "   - Test API:         curl http://localhost:8085"
echo "   - Database:         sudo -u postgres psql $DB_NAME"
echo ""
echo "📝 To update code:"
echo "   cd $APP_DIR"
echo "   git pull origin main"
echo "   go build -o mockhu-api cmd/api/main.go"
echo "   sudo systemctl restart mockhu"
echo ""
