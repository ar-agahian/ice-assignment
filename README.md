# ice-assignment

## Setup

1. Copy the `.env.example` file to the `/cmd/server` directory:
   ```bash
   cp .env.example cmd/server/.env
   ```

2. Run dependencies and start the application:
   ```bash
   make run
   ```

   This will:
   - Start all required services (MySQL, Redis, S3) using Docker Compose
   - Run the application server

3. Stop the application and dependencies:
   ```bash
   make stop
   ```

### Generate Mocks
```bash
make mocks
```

### Run Tests
```bash
make test
```

### Run Benchmark Tests
```bash
make benchmark
```

## API Endpoints

### 1. Upload File
**POST** `/api/asset`

Upload a file to S3 storage.

**Request:**
- Form field: `file` (the file to upload)

**Response:**
```json
{
  "fileId": "uuid-string"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/asset \
  -F "file=@/path/to/file.pdf"
```

### 2. Create Todo
**POST** `/api/todo`

Create a new todo item.

**Request:**
```json
{
  "description": "Complete the assignment",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-uuid-string"
}
```

**Response:**
```json
{
  "id": "uuid-string",
  "description": "Complete the assignment",
  "dueDate": "2024-12-31T23:59:59Z",
  "fileId": "optional-uuid-string"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/todo \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Complete the assignment",
    "dueDate": "2024-12-31T23:59:59Z"
  }'
```
