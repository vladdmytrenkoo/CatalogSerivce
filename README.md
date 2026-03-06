# Product Catalog Service

This service implements product catalog management according to Clean Architecture and DDD principles. It provides a robust gRPC API for managing products, including lifecycle management (creation, activation, archiving) and a flexible discount system.

## Technology Stack
- **Go 1.26**
- **gRPC** for high-performance API
- **Google Cloud Spanner** (using official Emulator)
- **Viper** for hierarchical configuration (YAML, Env)
- **Committer Pattern** for atomic transactions (Golden Mutation)
- **Transactional Outbox** for reliable event delivery

## Quick Start

### 1. Setup Infrastructure
Start the Spanner emulator using Makefile:
```bash
make docker-up
```

### 2. Run Application
Migrations make automatic whe app starts.
Migrations are performed automatically on every server start.
```bash
make run
```

### 3. Run Tests
The project includes domain unit tests and complex E2E workflow tests:
```bash
make test
```

---

## gRPC API Reference

The service is accessible at `localhost:50051`. It supports **gRPC Reflection**, so tools like Postman or BlumRPC will automatically discover all methods.

### Product Lifecycle

#### `CreateProduct`
Creates a new product in `DRAFT` status.
- **Request Body:**
```json
{
  "name": "IPhone 15 Pro",
  "description": "Powerful smartphone",
  "category": "Electronics",
  "base_price_numerator": 1200,
  "base_price_denominator": 1
}
```

#### `ActivateProduct`
Makes the product `ACTIVE`. Only active products are visible in `ListProducts`.
- **Request Body:**
```json
{
  "product_id": "UUID-FROM-CREATE"
}
```

#### `DeactivateProduct`
Moves product to `ARCHIVED` status.
- **Request Body:**
```json
{
  "product_id": "UUID"
}
```

### Management & Discounts

#### `UpdateProduct`
Updates product details. Empty fields are ignored (previous values are kept).
- **Request Body:**
```json
{
  "product_id": "UUID",
  "name": "New Name",
  "category": "New Category"
}
```

#### `ApplyDiscount`
Applies a percentage discount for a specific time range.
- **Request Body:**
```json
{
  "product_id": "UUID",
  "discount_percent": "15.5",
  "start_date": "2026-03-06T00:00:00Z",
  "end_date": "2026-12-31T23:59:59Z"
}
```

#### `RemoveDiscount`
Clears any active discount from the product.

### Reading Data

#### `GetProduct`
Returns full product details including calculated `effective_price`.
- **Request Body:**
```json
{
  "product_id": "UUID"
}
```
#### `ListProducts`
Returns a paginated list of **Active** products.
- **Request Body:**
```json
{
  "category": "Electronics",
  "page_size": 10
}
```
