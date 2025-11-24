# Todo List Backend API

## Project Overview

Backend API untuk aplikasi Todo List yang dibangun dengan Go (Golang) menggunakan framework Gin dan database PostgreSQL. API ini menyediakan endpoint RESTful untuk mengelola todos dan categories dengan fitur pagination, search, dan filtering.

## Features Implemented

### Todo Management
- Create, Read, Update, Delete (CRUD) todos
- Toggle todo completion status
- Pagination untuk list todos
- Search todos berdasarkan title
- Filter todos berdasarkan category, priority, dan completion status
- Sorting todos
- Validasi mandatory fields (title, category_id, priority)

### Category Management
- Create, Read, Update, Delete (CRUD) categories
- Validasi unique category name
- Prevent delete category jika masih digunakan oleh todos
- Default color untuk category

### API Features
- Standardized API response format (code, status, message, data)
- Error handling yang konsisten
- CORS middleware
- Health check endpoint

## Setup dan Installation

### Prerequisites

1. **Go** (versi 1.23 atau lebih tinggi)
   ```bash
   go version
   ```

2. **PostgreSQL** (versi 12 atau lebih tinggi)
   ```bash
   psql --version
   ```

3. **Git** (untuk clone repository)

### Step-by-Step Installation

#### 1. Clone Repository

```bash
git clone <repository-url>
cd todolist/be
```

#### 2. Install Dependencies

```bash
go mod download
```

#### 3. Setup Database

Buat database PostgreSQL:

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE todolist_db;

# Keluar dari psql
\q
```

#### 4. Konfigurasi Environment Variables

Buat file `.env` di folder `be/`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todolist_db
DB_SSLMODE=disable
PORT=8080
```

**Catatan:** Jika menggunakan default PostgreSQL (user: postgres, password: postgres), file `.env` tidak wajib.

#### 5. Run Database Migration

**PENTING:** Migration harus dijalankan sebelum start server!

```bash
go run cmd/migrate/main.go -action=up
```

Output yang diharapkan:
```
Database connected successfully
Database migration completed successfully
Migration up completed successfully
```

## How to Run Application Locally

### Cara 1: Menggunakan Go Run (Recommended untuk Development)

```bash
cd be
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

### Cara 2: Build dan Run Binary

```bash
cd be

# Build binary
go build -o server cmd/server/main.go

# Run binary
./server
```

### Cara 3: Menggunakan Docker (Optional)

```bash
cd be

# Build dan start containers
sudo docker-compose up --build -d

# Run migration
sudo docker-compose exec backend ./migrate -action=up

# Check status
sudo docker-compose ps
```

### Verify Server Running

Test health check endpoint:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "code": 200,
  "status": "success",
  "message": "Todo List API is running"
}
```

## How to Run Tests

Tests belum diimplementasikan. Untuk menambahkan tests:

```bash
# Run all tests
go test ./...

# Run tests dengan coverage
go test -cover ./...

# Run tests dengan verbose output
go test -v ./...
```

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Response Format

#### Success Response
```json
{
  "code": 200,
  "status": "success",
  "message": "Success message",
  "data": {}
}
```

#### Error Response
```json
{
  "code": 400,
  "status": "error",
  "message": "Error message"
}
```

#### Paginated Response
```json
{
  "code": 200,
  "status": "success",
  "message": "Success message",
  "data": [],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Endpoints

#### Health Check
```
GET /health
```

#### Todos

**Get All Todos**
```
GET /api/todos
Query Parameters:
  - page (int, default: 1)
  - limit (int, default: 10)
  - search (string, optional)
  - sort_by (string, optional)
  - sort_order (asc|desc, optional)
  - category_id (int, optional)
  - priority (high|medium|low, optional)
  - completed (bool, optional)
```

**Get Todo by ID**
```
GET /api/todos/:id
```

**Create Todo**
```
POST /api/todos
Body:
{
  "title": "string (required)",
  "description": "string (optional)",
  "category_id": "number (required)",
  "priority": "high|medium|low (required)",
  "due_date": "ISO 8601 string (optional)"
}
```

**Update Todo**
```
PUT /api/todos/:id
Body:
{
  "title": "string (optional)",
  "description": "string (optional)",
  "category_id": "number (optional)",
  "priority": "high|medium|low (optional)",
  "completed": "bool (optional)",
  "due_date": "ISO 8601 string (optional)"
}
```

**Delete Todo**
```
DELETE /api/todos/:id
```

**Toggle Todo Complete**
```
PATCH /api/todos/:id/complete
```

#### Categories

**Get All Categories**
```
GET /api/categories
```

**Get Category by ID**
```
GET /api/categories/:id
```

**Create Category**
```
POST /api/categories
Body:
{
  "name": "string (required)",
  "color": "hex color string (optional, default: #3B82F6)"
}
```

**Update Category**
```
PUT /api/categories/:id
Body:
{
  "name": "string (optional)",
  "color": "hex color string (optional)"
}
```

**Delete Category**
```
DELETE /api/categories/:id
```

### Example API Calls

```bash
# Get all todos
curl http://localhost:8080/api/todos

# Create todo
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete coding challenge",
    "category_id": 1,
    "priority": "high"
  }'

# Get all categories
curl http://localhost:8080/api/categories
```

## Technical Questions

### 1. Database Design

**Q: Bagaimana struktur database yang digunakan?**

A: Database menggunakan PostgreSQL dengan 2 tabel utama:

- **todos**: Menyimpan data todo dengan fields:
  - `id` (primary key)
  - `title` (required)
  - `description` (optional)
  - `completed` (boolean, default: false)
  - `category_id` (required, foreign key ke categories)
  - `priority` (enum: high, medium, low)
  - `due_date` (optional, timestamp)
  - `created_at`, `updated_at`, `deleted_at` (soft delete)

- **categories**: Menyimpan data category dengan fields:
  - `id` (primary key)
  - `name` (required, unique)
  - `color` (hex color string, default: #3B82F6)
  - `created_at`, `updated_at`

**Relationship:**
- Todos memiliki foreign key ke categories (many-to-one)
- Category tidak bisa dihapus jika masih digunakan oleh todos (RESTRICT)

### 2. API Design

**Q: Bagaimana API response format yang digunakan?**

A: API menggunakan standardized response format dengan struktur:

```json
{
  "code": 200,
  "status": "success|error",
  "message": "Human readable message",
  "data": {} // untuk success response
}
```

Untuk paginated response, ditambahkan field `pagination`:
```json
{
  "code": 200,
  "status": "success",
  "message": "Success message",
  "data": [],
  "pagination": {
    "current_page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

**HTTP Status Codes:**
- 200: Success
- 201: Created
- 400: Bad Request
- 404: Not Found
- 500: Internal Server Error

### 3. Error Handling

**Q: Bagaimana error handling diimplementasikan?**

A: Error handling dilakukan di beberapa layer:

1. **Validation Layer**: Menggunakan Gin binding untuk validasi request
2. **Service Layer**: Validasi business logic (category exists, priority valid, dll)
3. **Handler Layer**: Menggunakan utility functions untuk standardized error response

Semua error dikembalikan dalam format yang konsisten dengan code dan message yang jelas.

### 4. Database Migration

**Q: Bagaimana database migration dihandle?**

A: Database migration dilakukan secara manual menggunakan command terpisah:

```bash
go run cmd/migrate/main.go -action=up   # untuk migrate
go run cmd/migrate/main.go -action=down # untuk rollback
```

Migration files disimpan di folder `migrations/up/` dan `migrations/down/` dengan format SQL.

**Keuntungan pendekatan ini:**
- Kontrol penuh atas kapan migration dijalankan
- Bisa review SQL migration sebelum dijalankan
- Lebih mudah untuk rollback jika ada masalah

## Project Structure

```
be/
├── cmd/
│   ├── migrate/        # Migration command
│   └── server/         # Server entry point
├── internal/
│   ├── config/         # Configuration
│   ├── database/       # Database connection
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Middleware (CORS)
│   ├── models/         # Data models & DTOs
│   ├── router/         # Route setup
│   └── services/       # Business logic
├── migrations/         # SQL migration files
│   ├── up/            # Migration up
│   └── down/          # Migration down
├── pkg/
│   └── utils/          # Utility functions (response)
├── .env               # Environment variables (optional)
├── docker-compose.yml # Docker setup (optional)
├── Dockerfile         # Docker image (optional)
└── go.mod            # Go dependencies
```

## Troubleshooting

### Database Connection Error
- Pastikan PostgreSQL sudah running
- Cek kredensial di `.env` atau environment variables
- Pastikan database `todolist_db` sudah dibuat

### Port Already in Use
- Ubah PORT di `.env` atau environment variable
- Atau kill process yang menggunakan port 8080:
  ```bash
  lsof -ti:8080 | xargs kill -9
  ```

### Migration Error
- Pastikan migration dijalankan sebelum start server
- Cek apakah migration files sudah benar
- Pastikan database connection sudah benar

## License

MIT License

