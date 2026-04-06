# GoHeadless CMS

A Headless CMS implementation built with Go Fiber and MongoDB, following Clean Architecture and SOLID principles.

## Features
- **Dynamic Collection Management**: Systematically define data structures (collections, fields, and access control policies).
  - **Editable Schemas**: Safely mutate collection definitions at any point without losing integrity (protects identifiers).
- **Advanced Field System**: 
  - Standard Inputs (Text, Numbers, Booleans, Relational bindings).
  - Chronological inputs including new strictly-validated `daterange` fields.
  - Recursive layouts via structural fields (Groups, Repeaters, Tabs, Grid).
- **Media & File Handling**: Dedicated file management API with on-the-fly Image-to-WebP compression.
- **Dynamic Content Engine**: Seamless CRUD operations directly matching the schemas you've outlined. Includes nested data preservation and validation.
- **Modern Admin Panel**: A high-performance, reload-free SvelteKit SPA bridging you directly to your content and schemas.
- **SOLID Blueprint**:
    - **Single Responsibility (SRP)**: Handlers, Services, and Repositories are strictly segregated.
    - **Open/Closed (OCP)**: Field validation blocks can extend cleanly over time.
    - **Interface Segregation (ISP)**: Split abstractions like `CollectionWriter` vs `CollectionReader`.
    - **Dependency Inversion (DIP)**: Top-level domain structs hold no dependencies on physical drivers (MongoDB).

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
- **Admin Panel**: http://localhost:3030/
- **API**: http://localhost:3030/api/v1/
- **Swagger Docs**: http://localhost:3030/docs/

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

