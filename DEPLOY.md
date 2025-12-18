# Deployment Guide

This guide assumes you are deploying to an Ubuntu server with Nginx installed.

## 1. Build the Application

Run the following commands on your local machine (or the server if you have the tools installed) to build the binaries.

### Build Frontend
```bash
cd client
npm install
npm run build
# This creates a 'dist' folder with the static files
```

### Build Backend
```bash
cd server
# Build for Linux (if you are on Mac/Windows)
GOOS=linux GOARCH=amd64 go build -o price-is-right-server main.go
# If you are already on Linux, just run:
# go build -o price-is-right-server main.go
```

## 2. Prepare Directories on Server

SSH into your server and create the necessary directories.

```bash
# Directory for the backend executable
sudo mkdir -p /opt/price-is-right

# Directory for the frontend static files
sudo mkdir -p /var/www/price-is-right
```

## 3. Upload Files

Copy the files from your machine to the server.

1.  Copy the **backend binary** (`server/price-is-right-server`) to `/opt/price-is-right/server`.
2.  Copy the **frontend build** (`client/dist/*`) to `/var/www/price-is-right/`.

## 4. Setup Systemd Service

1.  Copy the service file:
    ```bash
    sudo cp deployment/price-is-right.service /etc/systemd/system/
    ```
2.  Reload systemd and start the service:
    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable price-is-right
    sudo systemctl start price-is-right
    ```
3.  Check status:
    ```bash
    sudo systemctl status price-is-right
    ```

## 5. Configure Nginx

1.  Copy the nginx config:
    ```bash
    sudo cp deployment/nginx.conf /etc/nginx/sites-available/price-is-right
    ```
2.  Edit the config to set your domain:
    ```bash
    sudo nano /etc/nginx/sites-available/price-is-right
    # Change 'server_name your-domain.com' to your actual domain
    ```
3.  Enable the site:
    ```bash
    sudo ln -s /etc/nginx/sites-available/price-is-right /etc/nginx/sites-enabled/
    ```
4.  Test and reload Nginx:
    ```bash
    sudo nginx -t
    sudo systemctl reload nginx
    ```

## 6. Done!

Visit your domain in the browser. The frontend should load, and it should successfully connect to the backend via the `/ws` path.
