# Locallyn Backend

> A location-aware community intelligence backend for surfacing nearby incidents, broadcasts, and trusted local updates in real time.

Locallyn is not built like a typical ecommerce CRUD app. The backend is designed around a more interesting problem: **how do you turn location-based user reports into a useful, searchable, trust-aware local feed?**

It combines geospatial storage, incident grouping, ranking logic, authenticated user identity, Redis-backed verification, and Elasticsearch-powered discovery into a clean Go service.

---

## Why This Project Stands Out

- **Geospatial-first design** using PostGIS geography points, distance filters, and spatial indexes.
- **Incident clustering** that finds or creates nearby incidents by category instead of treating every report as an isolated post.
- **Trust-aware feed ranking** based on post trust, recency, confirmations, feedback, reports, and proximity.
- **Searchable local feed** using Elasticsearch multi-match search with geo-distance filters and cursor-based `search_after` pagination.
- **Transactional consistency** for workflows that update multiple entities, such as post creation, incident counters, confirmations, feedback, and trust scores.
- **Layered Go architecture** with route, handler, service, repository, DTO, and infrastructure packages separated by domain.
- **Production-oriented infrastructure choices** with PostgreSQL/PostGIS, Redis, Elasticsearch, JWT authentication, Docker Compose, and environment-driven config.

---

## Core Idea

Locallyn lets users publish two kinds of local updates:

1. **Incident posts**: reports tied to a real-world event such as safety alerts, road issues, disruptions, or neighborhood activity.
2. **Broadcast posts**: general local updates visible to nearby users without being attached to an incident.

When a user creates an incident post, the backend checks whether a nearby incident already exists for the same category. If yes, the post joins that incident. If not, a new incident is created. This creates a feed that feels organized around real local events instead of becoming a noisy list of disconnected posts.

---

## Tech Stack

| Area | Technology | Purpose |
| --- | --- | --- |
| Language | Go | Backend service implementation |
| HTTP Framework | Gin | Routing, middleware, JSON APIs |
| Database | PostgreSQL + PostGIS | Relational data and geospatial queries |
| Cache | Redis | Verification token storage and short-lived auth support |
| Search | Elasticsearch | Full-text and geo-filtered feed search |
| Auth | JWT + bcrypt | Stateless access tokens and secure password hashing |
| DB Driver | pgx | PostgreSQL access from Go |
| IDs | UUID / NanoID | Stable resource identifiers and verification tokens |
| Infra | Docker Compose | Local Postgres, Redis, and Elasticsearch setup |
| Config | godotenv | Environment-based runtime configuration |

---

## Architecture

```text
cmd/server
   |
   v
api routes
   |
   v
domain handlers
   |
   v
domain services
   |
   v
repositories + infrastructure clients
   |
   +--> PostgreSQL/PostGIS
   +--> Redis
   +--> Elasticsearch
```

The codebase is organized by business domain:

```text
internal/
  user/           authentication, profile, verification, trust profile
  post/           post creation, fetching, feed document indexing
  incident/       nearby incident lookup and incident creation
  feed/           ranked location feed, incident posts, search feed
  interactions/   confirmations, feedback, reports, trust updates
  common/         JWT auth, constants, shared utilities

pkg/
  database/       PostgreSQL connection, schema, transactions
  redisConfig/    Redis client bootstrap
  elasticsearch/  index mapping, document building, indexing, search
```

---

## Key Backend Flows

### 1. Incident-Aware Post Creation

When an authenticated user creates an incident post:

- The service validates the JWT claims.
- A nearby incident is searched using latitude, longitude, radius, and category.
- The post is created inside a database transaction.
- Incident and user post counters are updated.
- A denormalized feed document is indexed into Elasticsearch.

This keeps the relational source of truth consistent while making the feed searchable.

### 2. Location-Based Feed Ranking

The feed is not just sorted by time. Nearby incidents are scored using:

- Trust score
- Freshness
- Post and confirmation activity
- Distance from the requesting user

Incident posts also use ranking signals such as:

- Post trust score
- Recency
- Feedback/report activity
- Distance from the incident center

### 3. Search With Geo Constraints

Feed search uses Elasticsearch to combine:

- Full-text matching on post content
- Boosted incident title matching
- Incident category matching
- Geo-distance filtering
- Deleted/flagged/expired post filters
- Stable cursor pagination with `search_after`

### 4. Trust and Moderation Signals

Users can:

- Confirm incidents
- Mark posts as helpful or misleading
- Report posts

The backend prevents users from giving feedback or reporting their own posts, and updates trust-related aggregates transactionally.

---

## API Capabilities

The API is grouped around five backend capabilities rather than simple entity CRUD:

- **Authentication and profiles**: signup, Redis-backed verification, login, JWT-protected profile creation and updates.
- **Local post publishing**: create incident or broadcast posts with location, identity mode, media, expiry, and trust metadata.
- **Incident intelligence**: discover nearby incidents, group related posts, track confirmations, and maintain incident-level counters.
- **Feed discovery**: serve ranked nearby feeds, fetch paginated incident conversations, and search local posts with geo-aware relevance.
- **Community signals**: collect confirmations, helpful/misleading feedback, and reports while preventing self-feedback abuse.

---

## Data Model Highlights

- `users`: authentication identity and verification status
- `user_profiles`: public profile, trust score, activity counters
- `incidents`: geospatial event clusters with confirmation and trust metadata
- `posts`: incident or broadcast updates with identity mode and expiry
- `post_feedback`: helpful/misleading feedback with one vote per user per post
- `incident_confirmations`: one confirmation per user per incident
- `post_reports`: moderation reports

PostGIS indexes are used on incident and post locations for efficient nearby queries.

---

## Local Development

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### Start dependencies

```bash
docker compose up -d
```

This starts:

- PostgreSQL + PostGIS on host port `5434`
- Redis on host port `6381`
- Elasticsearch on host port `9200`

### Environment example

Create a local `.env` file using this shape:

```env
PORT=<server-port>
POSTGRES_URL=<postgres-connection-url>
REDIS_ADDR=<redis-host-and-port>
ELASTICSEARCH_URL=<elasticsearch-url>
FEED_SEARCH_INDEX=<feed-index-name>
JWT_SECRET=<jwt-signing-secret>
ACCESS_TOKEN_EXPIRY=<access-token-expiry-in-minutes>
```

### Run the server

```bash
go run ./cmd/server
```

Health check:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

---

## Engineering Notes

- The app boots PostgreSQL, Redis, and Elasticsearch clients before starting the Gin server.
- Elasticsearch feed index bootstrapping happens during server startup.
- Database writes that affect multiple aggregates use transaction boundaries.
- Cursor pagination is used instead of offset pagination for feed and search stability.
- Authenticated routes are protected by JWT middleware.
- Passwords are stored using bcrypt, never as plain text.

---

## What This Demonstrates

This project demonstrates backend engineering beyond basic CRUD:

- Designing around geospatial product requirements
- Modeling event-based social data
- Building ranked feeds with explicit scoring signals
- Combining relational consistency with search indexing
- Structuring a Go service around clean domain boundaries
- Using Redis, Elasticsearch, PostGIS, JWT, and Docker in one cohesive backend

---

## Future Enhancements

- Background jobs for re-indexing feed documents
- Refresh-token based auth flow
- Rate limiting for post creation and reporting
- Outbox pattern for reliable Elasticsearch indexing
- Admin moderation APIs
- Automated integration tests with containerized dependencies
