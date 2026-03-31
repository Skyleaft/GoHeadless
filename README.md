# GoHeadless CMS

A Headless CMS implementation built with Go Fiber and MongoDB, following Clean Architecture and SOLID principles.

## Features
- **Dynamic Collection Management (Metadata)**: Define your own data structures (collections and fields).
- **Dynamic Content Engine**: CRUD operations on the collections you've defined, with automatic schema validation.
- **Admin Panel**: SvelteKit SPA for managing collections and content.
- **SOLID Principles**:
    - **Single Responsibility (SRP)**: Each component has a single purpose.
    - **Open/Closed (OCP)**: The system supports multiple field types and can be extended without altering the core.
    - **Interface Segregation (ISP)**: Split repository interfaces (e.g., Reader/Writer).
    - **Dependency Inversion (DIP)**: High-level business logic doesn't depend on low-level database details.

## Tech Stack
- **Go 1.21+**
- **Fiber v3** (Web Framework)
- **MongoDB** (Storage)
- **Go Driver for MongoDB**
- **SvelteKit** (Admin Panel)
- **Nginx** (Reverse Proxy)
- **Docker & Docker Compose** (Containerization)

## Directory Structure
```
/cmd/api            # Fiber app initialization
/internal
  /domain           # Domain entities & repository interfaces
  /collection       # Collection Manager (Service, Repo, Handler)
  /content          # Content Engine (Service, Repo, Handler)
  /platform         # MongoDB connectivity
/Panel              # Admin panel (SvelteKit SPA)
/nginx              # Nginx configuration
```

## Development Setup

### Prerequisites
- Go 1.21+
- Node.js 20+
- MongoDB

### Running the API
1. Ensure MongoDB is running (Default: `mongodb://localhost:27017`).
2. Set environment variables (Optional):
   - `MONGO_URI`: MongoDB connection string.
   - `DB_NAME`: Database name (default: `goheadless_cms`).
   - `PORT`: Server port (default: `3000`).
3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### Running the Admin Panel
1. Navigate to the Panel directory:
   ```bash
   cd Panel
   ```
2. Install dependencies:
   ```bash
   pnpm install
   ```
3. Start the development server:
   ```bash
   pnpm dev
   ```

## Docker Deployment

The project includes complete Docker configuration for production deployment.

### Quick Start
```bash
docker-compose up --build -d
```

### Services
| Service | Description | Port |
|---------|-------------|------|
| **MongoDB** | Database | 27017 |
| **API** | Go backend (internal) | - |
| **Nginx** | Reverse proxy + admin panel | 80 |

### Access Points
- **Admin Panel**: http://localhost/admin
- **API**: http://localhost/api/v1
- **Swagger Docs**: http://localhost/docs

### Environment Variables
Create a `.env` file:
```env
MONGO_ROOT_USERNAME=admin
MONGO_ROOT_PASSWORD=your_secure_password
DB_NAME=goheadless_cms
APP_PORT=80
```

### Docker Commands
```bash
# View logs
docker-compose logs -f

# Rebuild and restart
docker-compose up --build -d

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

See [DOCKER-README.md](DOCKER-README.md) for detailed documentation.

