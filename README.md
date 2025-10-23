# GoCatalog

GoCatalog is a RESTful API service built with `GoLang` and the `Gin framework`.

# Table of Contents

1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Getting Started (Installation)](#getting-started-installation)
4. [API Endpoints](#api-endpoints)
   1. [Health Check](#health-check)
   2. [Upload CSV](#upload-csv)
   3. [Get All Products](#get-all-products)
   4. [Filter Products](#filter-products)
5. [Optimizations](#optimizations)
   1. [Database-Level Optimizations](#database-level-optimizations)
   2. [Docker Optimizations](#docker-optimizations)
   3. [Rate Limiting](#rate-limiting)
6. [Testing](#testing)
7. [Future Implementation](#future-implementation)

### Features

- **REST API**: A simple API to interact with the product catalog.
- **CSV Upload**: Bulk import products into the database via a CSV file upload.
- **PostgreSQL Database**: Uses a PostgreSQL database for persistent storage.
- **Containerized**: Docker and Docker Compose for a consistent and isolated development environment.
- **Integration Tests**: Full test coverage to verify API endpoints and database interactions.
- **Pagination**: Supports `page` and `limit` query parameters for fetching large datasets efficiently.
- **Rate Limiting**: Token Bucket algorithm per IP to prevent abuse and ensure fair usage.

## Prerequisites

- [Docker (v28.5.1)](https://docs.docker.com/engine/release-notes/28/#2851)
- [Docker Compose (v2.40.0)](https://github.com/docker/compose/releases/tag/v2.40.0)

## Getting Started (Installation)

Follow these steps to get the project up and running on your local machine.

1. **Clone the repository**

   ```bash
   git clone https://github.com/dis70rt/GoCatalog
   cd GoCatalog
   ```
2. **Configure Environment**

   - Copy `.env.example` to `.env` and update credentials to your own.
     ```bash
     cp .env.example .env
     ```
3. **Configure Docker Compose**

   - If using a separate Docker Compose for your environment, copy `docker-compose.mock.yml` or `docker-compose.yml` and update it with your credentials if needed.
4. **Build and Run with Docker Compose**
   This command will build the Go application image and start the `gocatalog` and `postgres` services in detached mode.

   ```bash
   docker compose up --build -d
   ```

   The API will be running and accessible at `http://localhost:8080`.
5. **Stopping the services**
   To stop the running containers, use:

   ```bash
   docker compose down
   ```

You can manage services using the Makefile: `make up` (start), `make down` (stop), `make build` (build images)

## API Endpoints

<details>
  <summary>Health Check <a name="upload-csv"></summary>

#### [**Endpoint**: `GET /health`]()

[  **Description**: Checks the health of the API service.]()

#### **Endpoint**: `POST /upload`

  **Description**: Uploads a CSV file to bulk-insert products into the database. The CSV file must contain a header row.

- **Curl Example**
  ```bash
  curl -X POST http://localhost:8080/upload \
      -F "file=@/path/to/your/products.csv" \
      -H "Content-Type: multipart/form-data"
  ```
- **Success Response**:
  ```json
  {
      "stored": 150
  }
  ```
- **Failed Response**:
  ```json
    {
      "failed":[
          {
            "sku":"TSHIRT-RED-001",
            "name":"Classic Cotton T-Shirt",
            "brand":"StreamThreads",
            "size":"M",
            "color":"Red",
            "mrp":799,
            "price":499,
            "quantity":-5,
            "failed_reason":"Quantity invalid"
          }
      ],
      "stored":0
    }
  ```

</details>

<details>
  <summary>Get All Products <a name="get-all-products"></summary>

#### [**Endpoint**: `GET /products`]()

[  **Description**: Returns all stored products with pagination support. Use `page` and `limit` query parameters to control pagination.]()

- [**Query Parameters**:]()

  - [`page` (optional, default `1`) – Page number]()
  - [`limit` (optional, default `10`) – Number of products per page]()
- [**Curl Example**]()
- [**Success Response**:]()

#### [**Endpoint** `GET /products/search`]()

#### [**Description**]()

[Returns products filtered by brand, color, and price range, with pagination support. All filters are optional.]()

#### [**Query Parameters**]()

- [`brand` (optional) – Filter by product brand]()
- [`color` (optional) – Filter by product color]()
- [`minPrice` (optional) – Minimum price]()
- [`maxPrice` (optional) – Maximum price]()
- [`page` (optional, default `1`) – Page number]()
- [`limit` (optional, default `10`) – Number of products per page]()

##### [**Curl Example**]()

##### [**Success Response**]()

## [Optimizations]()

[GoCatalog has been built with performance and efficiency in mind. The following optimizations have been applied:]()

##### [Database-Level Optimizations]()

- [Bulk inserts for CSV uploads `(faster writes)`.]()
- [Frequently queried columns (sku, brand, color, price) are indexed to speed up searches.]()
- [Pagination and filtering done in SQL `(minimal memory usage)`]()

##### [Docker Optimizations]()

- [Multi-stage builds `(small image size)`]()
- [Alpine base `(lightweight containers)`]()

##### [Rate Limiting]()

- [Token Bucket algorithm for request limiting]()
- [Per-IP limits]()
- [Implemented as Gin middleware for all endpoints]()

### [Testing]()

[From the project root, run:]()

[This will:]()

- [Execute all tests in the repository]()
- [Show verbose output, including which tests passed or failed]()

### [Future Implementation]()

- [**Improved Project Structure**]()

[This makes the project easier to navigate, maintain, and scale as it grows.]()

- [**API Key Authorization**
  Implement authentication and authorization using API keys to secure endpoints and control access per client.]()
- [**Advanced Rate Limiting**
  Integrate with Redis for distributed environments.]()
