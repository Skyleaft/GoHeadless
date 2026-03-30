# GoHeadless CMS

A Headless CMS implementation built with Go Fiber and MongoDB, following Clean Architecture and SOLID principles.

## Features
- **Dynamic Collection Management (Metadata)**: Define your own data structures (collections and fields).
- **Dynamic Content Engine**: CRUD operations on the collections you've defined, with automatic schema validation.
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

## Directory Structure
```
/cmd/api            # Fiber app initialization
/internal
  /domain           # Domain entities & repository interfaces
  /collection       # Collection Manager (Service, Repo, Handler)
  /content          # Content Engine (Service, Repo, Handler)
  /platform         # MongoDB connectivity
```

## How to Run
1. Ensure MongoDB is running (Default: `mongodb://localhost:27017`).
2. Set environment variables (Optional):
   - `MONGO_URI`: MongoDB connection string.
   - `DB_NAME`: Database name (default: `goheadless_cms`).
   - `PORT`: Server port (default: `3000`).
3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## Example API Usage
### 1. Create a Collection
```bash
POST /api/v1/collections
{
  "name": "Products",
  "slug": "products",
  "fields": [
    { "name": "title", "type": "String", "required": true },
    { "name": "price", "type": "Float", "required": true }
  ]
}
```

### 2. Add Content to Collection
```bash
POST /api/v1/content/products
{
  "title": "Example Product",
  "price": 29.99
}
```
