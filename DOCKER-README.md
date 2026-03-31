# Docker Setup for GoHeadless CMS

This directory contains Docker configuration for running GoHeadless CMS with:
- **Go API** - Headless CMS backend
- **MongoDB** - Database
- **Nginx** - Reverse proxy and static file serving for the admin panel

## Architecture

```
Request → Nginx (:80) 
           ├── / → Static files (Svelte SPA)
           ├── /api/* → Go API (:3000)
           ├── /docs/* → Go API (Swagger)
           └── /uploads/* → Go API uploads
```

## Quick Start

### 1. Build and Start All Services

```bash
docker-compose up --build -d
```

### 2. Access the Application

- **Admin Panel**: http://localhost:3030/
- **API**: http://localhost:3030/api/v1/
- **Swagger Docs**: http://localhost:3030/docs/
- **MongoDB**: localhost:27017

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Database
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=your_secure_password
DB_NAME=goheadless_cms

# Ports
MONGO_PORT=27017
APP_PORT=3030         # Host port (external access)
NGINX_PORT=80        # Container internal port
```

## Docker Commands

### Start Services

```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d api
docker-compose up -d mongodb
docker-compose up -d nginx
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f mongodb
docker-compose logs -f nginx
```

### Rebuild

```bash
# Rebuild all services
docker-compose build

# Rebuild specific service
docker-compose build api
docker-compose build nginx
```

## File Structure

```
├── Dockerfile                 # Go API Dockerfile
├── .dockerignore              # Root dockerignore
├── docker-compose.yml         # Docker Compose configuration
├── nginx/
│   └── nginx.conf            # Nginx reverse proxy configuration
└── Panel/
    ├── Dockerfile            # Frontend development Dockerfile
    ├── Dockerfile.nginx      # Frontend + Nginx production Dockerfile
    └── .dockerignore         # Panel dockerignore
```

## Ports

| Service | Internal Port | External Port | Description |
|---------|--------------|---------------|-------------|
| Nginx   | 80           | 3030 (via APP_PORT) | Reverse proxy & static files |
| Go API  | 3000         | N/A (internal) | Backend API |
| MongoDB | 27017        | 27017 (via MONGO_PORT) | Database |

### Custom Port Configuration

To change the external port (e.g., 3030:80), set in your `.env`:
```env
APP_PORT=3030
NGINX_PORT=80
```
This maps host port 3030 to container port 80.

## Nginx Routing

| Path | Description |
|------|-------------|
| `/` | Admin panel SPA (frontend) |
| `/api/*` | Proxied to Go API |
| `/docs/*` | Proxied to Go API (Swagger) |
| `/uploads/*` | Proxied to Go API (uploaded files) |
