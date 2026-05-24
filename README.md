# Advanced Go API

This is the API documentation for the Advanced Go application.

## Base URL

All API endpoints are prefixed with `/api/v1`.

```
http://localhost:8080/api/v1
```

## Authentication

A valid JWT token must be provided in the `Authorization` header for all authenticated routes.

```
Authorization: Bearer <YOUR_JWT_TOKEN>
```

### Authentication Routes

#### Register a new user

- **POST** `/auth/register`
- **Description:** Creates a new user and returns a JWT token.
- **Roles:** None

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d 
'{'
  "email": "test@example.com",
  "password": "password123",
  "full_name": "Test User",
  "address": "123 Test St"
}'
```

**Sample Response:**

```json
{
  "access_token": "your.jwt.token"
}
```

#### Login

- **POST** `/auth/login`
- **Description:** Authenticates a user and returns a JWT token.
- **Roles:** None

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d 
'{'
  "email": "test@example.com",
  "password": "password123"
}'
```

**Sample Response:**

```json
{
  "access_token": "your.jwt.token"
}
```

#### Get current user

- **GET** `/auth/me`
- **Description:** Returns the profile of the currently authenticated user.
- **Roles:** Any authenticated user.

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### User Routes

#### Get user profile

- **GET** `/users/:id/profile`
- **Description:** Returns the public profile of a user.
- **Roles:** None

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/users/1/profile
```

#### Get current user (alternative)

- **GET** `/users/me`
- **Description:** Returns the profile of the currently authenticated user.
- **Roles:** Any authenticated user.

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/users/me \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

#### Update current user

- **PATCH** `/users/me`
- **Description:** Updates the profile of the currently authenticated user.
- **Roles:** Any authenticated user.

**Curl Command:**

```bash
curl -X PATCH http://localhost:8080/api/v1/users/me \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "full_name": "New Name"
}'
```

### Product Routes

#### List products

- **GET** `/products`
- **Description:** Returns a paginated list of active products.
- **Roles:** None

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/products?page=1&size=10
```

#### Get a single product

- **GET** `/products/:product_id`
- **Description:** Returns a single product by its ID.
- **Roles:** None

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/products/1
```

#### Create a product

- **POST** `/products`
- **Description:** Creates a new product.
- **Roles:** `AGENT`

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/products \
-H "Authorization: Bearer <AGENT_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "title": "New Product",
  "description": "A great new product",
  "price": 99.99,
  "is_active": true
}'
```

#### Update a product

- **PATCH** `/products/:product_id`
- **Description:** Updates a product.
- **Roles:** `AGENT` (must be the owner of the product)

**Curl Command:**

```bash
curl -X PATCH http://localhost:8080/api/v1/products/1 \
-H "Authorization: Bearer <AGENT_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "price": 129.99
}'
```

#### Delete a product

- **DELETE** `/products/:product_id`
- **Description:** Deletes a product.
- **Roles:** `AGENT` (must be the owner of the product)

**Curl Command:**

```bash
curl -X DELETE http://localhost:8080/api/v1/products/1 \
-H "Authorization: Bearer <AGENT_JWT_TOKEN>"
```

#### List product reviews

- **GET** `/products/:product_id/reviews`
- **Description:** Returns all reviews for a product.
- **Roles:** None

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/products/1/reviews
```

#### Create a product review

- **POST** `/products/:product_id/reviews`
- **Description:** Creates a review for a product.
- **Roles:** `CONSUMER` (must have purchased the product)

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/products/1/reviews \
-H "Authorization: Bearer <CONSUMER_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "rating": 5,
  "comment": "This product is amazing!"
}'
```

### Order Routes

#### Create an order

- **POST** `/orders`
- **Description:** Creates a new order.
- **Roles:** `CONSUMER`

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/orders \
-H "Authorization: Bearer <CONSUMER_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    }
  ]
}'
```

#### List my orders

- **GET** `/orders`
- **Description:** Returns a paginated list of orders for the authenticated user.
- **Roles:** Any authenticated user.

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/orders?page=1&size=10 \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

#### Get a single order

- **GET** `/orders/:order_id`
- **Description:** Returns a single order by its ID.
- **Roles:** The user who created the order or the agent who owns the products in the order.

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/orders/1 \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

#### Update order status

- **PATCH** `/orders/:order_id/status`
- **Description:** Updates the status of an order.
- **Roles:** `AGENT` (must be the owner of the products in the order)

**Curl Command:**

```bash
curl -X PATCH http://localhost:8080/api/v1/orders/1/status \
-H "Authorization: Bearer <AGENT_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "status": "SHIPPED"
}'
```

### Message Routes

#### Send a message

- **POST** `/messages/orders/:order_id`
- **Description:** Sends a message related to an order.
- **Roles:** The user who created the order or the agent who owns the products in the order.

**Curl Command:**

```bash
curl -X POST http://localhost:8080/api/v1/messages/orders/1 \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
-H "Content-Type: application/json" \
-d 
'{'
  "content": "When will my order arrive?"
}'
```

#### List messages

- **GET** `/messages/orders/:order_id`
- **Description:** Returns all messages for an order.
- **Roles:** The user who created the order or the agent who owns the products in the order.

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/messages/orders/1 \
-H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### Admin Routes

#### Get statistics

- **GET** `/admin/stats`
- **Description:** Returns administrative statistics.
- **Roles:** `ADMIN`

**Curl Command:**

```bash
curl -X GET http://localhost:8080/api/v1/admin/stats \
-H "Authorization: Bearer <ADMIN_JWT_TOKEN>"
```

### Monitoring Routes

#### Metrics

- **GET** `/metrics`
- **Description:** Returns Prometheus metrics.
- **Roles:** None

**Curl Command:**

```bash
curl -X GET http://localhost:8080/metrics
```

#### Profiling

- **GET** `/debug/pprof/`
- **Description:** Provides pprof profiling data.
- **Roles:** None

**Curl Command:**

```bash
# See pprof documentation for usage
# Example: go tool pprof http://localhost:8080/debug/pprof/profile
```

* Use it cuz i gaurentee it will work.
* code and readme both by my glitched mind
