# Technical Challenge – Rate Limiter (Go)

This file describes the **original challenge statement** (functional and technical requirements).  
For details of the concrete implementation in this repository (how to run, configure and test), see the main `README.md`.

---

## Goal

Build a **Rate Limiter** as an **HTTP middleware** capable of:

- Limiting requests by **IP address** or by **access token** (`API_KEY`).
- Persisting counters and time windows in **Redis** (via Docker).
- Keeping an extensible architecture (using the **Strategy** pattern) so that the persistence backend can be swapped in the future.

---

## Business rules (limitation criteria)

### Limitation by IP

Restrict the maximum number of **requests per second** received from a single IP address.

### Limitation by token

Restrict the maximum number of **requests per second** based on a unique access token.

- **Header:** the token is sent as `API_KEY: <TOKEN>`.

### Precedence (golden rule)

Token configuration **must override** IP configuration.

**Example:** global limit per IP = 10 req/s; a specific token = 100 req/s → the system must apply **100 req/s** for requests that include this token.

---

## Behavior when blocked

When the limit is exceeded (by IP or token):

| Aspect | Behavior |
|--------|----------|
| **HTTP** | Immediate **429** status code. |
| **Body** | Exactly: `you have reached the maximum number of requests or actions allowed within a certain time frame` |
| **Block time** | The offending IP or token is blocked for a **configurable** period (e.g., 5 minutes). While blocked, new requests must be **rejected**. |

---

## Technical requirements and architecture

### Middleware

The rate limiter logic must be implemented as an **HTTP middleware** that wraps the web server.

### Persistence (Redis)

Counters and time control must be stored in **Redis** (brought up via Docker Compose).

### Design pattern: Strategy

- Implement a **Strategy** for the persistence layer.
- **Redis is mandatory** for this challenge, but the architecture must allow swapping the storage by simply changing the strategy implementation (without changing business rules).

### Decoupling

The **business rules** of the limiter must be **separated** from the **middleware logic**:

- The middleware only orchestrates (extracts IP/token, calls the limiter, decides 200/429).
- The limiter applies rules for limits, counting, and blocking.

### Configuration

Everything must be configured via **environment variables** and/or a **`.env` file in the root**, for example:

- Maximum number of requests per second (IP and/or token).
- Block time after exceeding the limit.
- Redis connection parameters.

---

## How to run the project (challenge requirement)

The evaluator must be able to bring up the application and tests using **only Docker / Docker Compose**.

- Bring up application + Redis from the repository root (where `docker-compose.yaml` and `Dockerfile` are):

```bash
docker compose up --build
```

- The application must listen on port **8080**.

- Automated tests (example):

```bash
docker compose run --rm app go test ./...
```

Or an equivalent target/command documented in the project.

---

## Expected tests

Include tests that demonstrate:

1. **Effectiveness** of the limiter (429 after exceeding the limit, temporary blocking, etc.).
2. **Precedence token > IP** (a token with a higher limit than the IP must behave according to the token limit).

---

## How to configure the limiter

1. Provide a `.env.example` file at the repository root to be copied to `.env`.
2. Configure, via environment variables:
   - per-second limits (IP and token),
   - block duration,
   - Redis connection parameters.
3. (Optional) If the solution supports different limits per specific token, document how to register the **token → limit** map.

---

## How to change the persistence strategy

1. Define an interface (Strategy) used by the limiter domain (e.g., increment counter, check block, TTL).
2. Implement the current **Redis** strategy.
3. For another backend, create a new struct that implements the same interface and inject it into the middleware/limiter via composition or a factory — **without changing the limiter’s business rules**.

---

## Expected deliverables

| Item | Description |
|------|-------------|
| **Source code** | Complete implementation of the rate limiter in Go. |
| **Dockerfile** | Builds the Go application. |
| **docker-compose.yaml** | Brings up the app on port **8080** and Redis. |
| **Project README** | Explains configuration, strategies, and how to run the implementation. |
| **Tests** | Cover the limiter and token > IP precedence. |

---

## Delivery rules

- **Single repository:** only this project.
- **Main branch:** all code must be on the `main` branch.
- **Execution:** run the project and tests using **Docker / Docker Compose** only.

