# Go Event Registration System

## Architecture Explanation

This system strictly follows a Layered (Clean) Architecture approach:

1.  **Models (`internal/models`)**: Database schema definitions using standard GORM tags and structs.
2.  **Repositories (`internal/repositories`)**: Abstraction layer for all PostgreSQL operations. All DB queries are encapsulated here keeping other layers unaware of DB drivers.
3.  **Services (`internal/services`)**: The core business logic layer. Validation, calculations, concurrency-handling, and complex multi-repository logic sit here.
4.  **Handlers (`internal/handlers`)**: The Presentation Layer. They translate HTTP requests from Gin into service calls, and map service results back to HTTP JSON responses.
5.  **Middleware & Router (`internal/middleware` & `internal/router`)**: Routing and security (JWT interception) layers. We secure different route groups using dedicated role-check middleware.

## ER Diagram (Text Representation)

*   **User**: `id (UUID, PK)`, `name`, `email (Unique)`, `password_hash`, `role (ENUM)`, `is_active`
    *   *1:N* with **Event** (Organizer)
    *   *1:N* with **Registration**
    *   *1:N* with **Waitlist**
*   **Event**: `id (UUID, PK)`, `title`, `description`, `event_date`, `capacity`, `seats_remaining`, `organizer_id (FK)`, `status`
    *   *1:N* with **Registration**
    *   *1:N* with **Waitlist**
*   **Registration**: `id (UUID, PK)`, `user_id (FK)`, `event_id (FK)`, `status (ENUM: CONFIRMED/CANCELLED)`
    *   Unique constraint on `(user_id, event_id)`
*   **Waitlist**: `id (UUID, PK)`, `user_id (FK)`, `event_id (FK)`, `position`
*   **AuditLog**: `id (UUID, PK)`, `actor_id (FK)`, `action`, `entity_type`, `entity_id`, `timestamp`

## Concurrency Strategy: `SELECT FOR UPDATE` vs Optimistic Locking

### 1. Pessimistic Locking (`SELECT FOR UPDATE`) - *Chosen Strategy*
When a user attempts to book an event, the system immediately begins a transaction and executes:
```sql
SELECT * FROM events WHERE id = ? FOR UPDATE;
```
**Why was this chosen?**
This locks the specific Event row exclusively for that transaction. Any concurrent requests trying to book the same event must wait until the current transaction commits or rolls back. This provides a strict, serial execution of bookings at the database level, completely preventing race conditions and negative seats. The code logic inside Go can safely read `seats_remaining`, decrement it, and save it without worrying that another thread changed the value in the meantime.

### 2. Alternative Approach: Optimistic Locking
Optimistic locking usually relies on a `version` integer column.
```sql
UPDATE events SET seats_remaining = seats_remaining - 1, version = version + 1 WHERE id = ? AND version = ?
```
If the version doesn't match (because another thread updated it), the update fails.
**Trade-offs vs Pessimistic:**
*   **Pros**: No database-level row locks are held. Performs better under low contention.
*   **Cons**: Under very high concurrent scale (e.g. 100 users booking 1 seat exactly at the same millisecond), optimistic locking fails 99 of the requests resulting in heavy RETRY logic, which taxes application resources and database IOPS. Pessimistic locking handles high-contention spikes perfectly.

## Scalability Considerations

1.  **Multiple App Instances**: Because the lock (`FOR UPDATE`) is managed by the PostgreSQL database engine, this approach is perfectly safe across horizontally scaled stateless application instances (e.g., Kubernetes pods running the Go app). Lock contention is solved at the Data Tier.
2.  **Trade-offs of Pessimistic Locking**: The biggest trade-off is latency during high contention. If 1,000 users hit one specific event simultaneously, their transaction requests queue up within Postgres. This could lead to temporary DB connection exhaustion if the connection pool isn't appropriately sized. 

## Concurrency Simulation (`POST /admin/events/:id/simulate?users=N`)
The application includes a testing handler acting as a stress-tester. When invoked, it:
1.  Creates N dummy users.
2.  Spawns N concurrent Goroutines.
3.  Each Goroutine invokes the `BookEvent` transaction concurrently.
4.  By waiting for the `WaitGroup`, it tallies successfully registered, waitlisted, and failed accounts, proving zero overbooking errors.
