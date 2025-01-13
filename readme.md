# Technical Documentation

## Project Overview

This project is a microservice-based system for processing, storing, and retrieving user data. It uses a combination of PostgreSQL, Redis, and RabbitMQ for storage and messaging, along with Go Fiber for API endpoints.

### Features

1. Upload CSV files and process user data.
2. Store user data in PostgreSQL for persistence.
3. Cache user data in Redis for fast retrieval.
4. Use RabbitMQ for message queuing to decouple data ingestion and processing.
5. API endpoints for retrieving user data with filtering and pagination.

---

## System Architecture

### Components

1. **API Service (`api.go`)**
   - Handles user requests to fetch data.
   - Queries Redis for cached user data.
   - Supports filtering by first name, last name, and email.
   - Implements pagination.

2. **Consumer Service (`consumer.go`)**
   - Consumes messages from RabbitMQ.
   - Parses and validates incoming user data.
   - Stores data in PostgreSQL and caches it in Redis.

3. **Database Service (`db.go`)**
   - Initializes and ensures the users table in PostgreSQL.
   - Connects to Redis for caching.

4. **Main Service (`main.go`)**
   - Handles CSV uploads.
   - Parses CSV and publishes data to RabbitMQ.

5. **Docker Compose**
   - Orchestrates the API, Consumer, PostgreSQL, Redis, and RabbitMQ services.

6. **Dockerfiles**
   - Contains Dockerfiles for building the API and Consumer services.

7. **Web**
   - head to localhost:8080 to get the updation of events.

---

## Prerequisites

### Software Requirements

- **Docker**: Install Docker and Docker Compose.
- **Go**: Install Go 1.18 or higher.

### Environment Variables

Ensure the following environment variables are set in the services:

- `POSTGRES_USER`: PostgreSQL username.
- `POSTGRES_PASSWORD`: PostgreSQL password.
- `POSTGRES_DB`: Database name.

---

## Installation and Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd <repository-directory>
```

### 2. Build and Start the Docker Containers

```bash
docker-compose up --build
```

### 3. Access Services

- API Service: [http://localhost:3001](http://localhost:3001)
- Producer Service: [http://localhost:3002](http://localhost:3002)
- Consumer Service: [http://localhost:3000](http://localhost:3000)
- Web Service: [http://localhost:8080](http://localhost:8080)
- PostgreSQL: [http://localhost:5432](http://localhost:5432)
- Redis: [http://localhost:6379](http://localhost:6379)
- RabbitMQ: [http://localhost:5672](http://localhost:5672)
- RabbitMQ Management UI: [http://localhost:15672](http://localhost:15672) (Username: `guest`, Password: `guest`)

---

## API Endpoints

- You can import and use [postman collection](https://github.com/ashvin-kumbhani/go-test/blob/main/Go-test-Backend.postman_collection.json) for detailed api documentation.

### 1. Upload CSV

**URL**: `/upload`
**Method**: `POST`
**Description**: Upload a CSV file for processing.
**Payload**:

- `file`: CSV file containing user data.

### 2. Fetch Users

**URL**: `/users`
**Method**: `GET`
**Description**: Retrieve user data with optional filters and pagination.
**Query Parameters**:

- `first_name` (optional): Filter by first name.
- `last_name` (optional): Filter by last name.
- `email` (optional): Filter by email.
- `page` (optional): Page number (default: 1).
- `page_size` (optional): Number of users per page (default: 10).

---

## Directory Structure

```
.
├── api.go              # API service implementation
├── consumer.go         # Consumer service for RabbitMQ
├── db.go               # Database and Redis initialization
├── main.go             # CSV upload and RabbitMQ producer
├── Dockerfile.api      # Dockerfile for API service
├── Dockerfile.consumer # Dockerfile for Consumer service
├── docker-compose.yml  # Docker Compose configuration
```

---

## Workflow

### Data Flow

1. User uploads a CSV file via `/upload` endpoint.
2. Records from the CSV are published to RabbitMQ.
3. The consumer service consumes messages from RabbitMQ:
   - Validates and processes the data.
   - Stores data in PostgreSQL.
   - Caches data in Redis.
4. API service retrieves data from Redis and returns it to the user.

### Queue Usage

RabbitMQ decouples the ingestion and processing of user data, ensuring scalability and fault tolerance.

---

## Running Locally

### Step 1: Start All Services

```bash
docker-compose up --build
```

### Step 2: Test Upload Endpoint

```bash
curl -F "file=@users.csv" http://localhost:3000/upload
```

### Step 3: Test Fetch Users Endpoint

```bash
curl "http://localhost:3001/users?page=1&page_size=10"
```

---

## Testing

### Unit Tests

Use Go's built-in testing tools to write and run tests for individual components.

### Integration Testing

1. Ensure all containers are running.
2. Test full workflow:
   - Upload a CSV file.
   - Verify data in PostgreSQL and Redis.
   - Query data using the `/users` endpoint.

---

## Improvements

1. Authentication and authorization for API endpoints.
2. I haven't use environment variables in the code fo ease of running it, I would use them to make the code more flexible for production.
3. I haven't encrypted the crucial data like email to make it searchable.
