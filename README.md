# Product Feeds

A TikTok-style vertical product feed app. Swipe through products full-screen, swipe horizontally through multiple product images, and get infinite scroll powered by a Go backend and Supabase.

## Stack

- **Backend:** Go + Gin + pgx (PostgreSQL driver)
- **Frontend:** Flutter
- **Database:** Supabase (PostgreSQL + Storage)

## Project Structure

```
product-feeds/
├── .env                  # Environment variables (see setup below)
├── backend/              # Go REST API
│   ├── main.go
│   ├── db/db.go          # Database connection pool
│   ├── handlers/         # HTTP handlers
│   └── models/           # Data models
└── frontend/             # Flutter app
    └── lib/
        ├── main.dart
        ├── models/
        ├── services/     # API client
        └── screens/      # UI
```

## Prerequisites

- Go 1.24+
- Flutter 3.11+
- A Supabase project with a `products` table

### Supabase `products` Table Schema

```sql
create table products (
  id           uuid primary key default gen_random_uuid(),
  name         text not null,
  price        numeric not null,
  currency     text not null,
  image_paths  jsonb not null,   -- array of image URLs
  source_url   text,
  company_name text,
  scraped_at   timestamp with time zone default now()
);
```

## Setup

### 1. Environment Variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgresql://<user>:<password>@<host>:<port>/postgres
SUPABASE_PROJECT_URL=https://<project-ref>.supabase.co
SUPABASE_PROJECT_API=<your-anon-or-service-key>
SUPABASE_DB_PASSWORD=<your-db-password>
SUPABASE_STORAGE_BUCKET=product-images
```

### 2. Backend

```bash
cd backend
go mod download
go run main.go
```

The API starts on `http://localhost:8080`.

**Endpoint:**

```
GET /api/products?limit=10&cursor=<uuid>
```

- `limit` — number of products to return (default 10, max 50)
- `cursor` — UUID of the last product seen, for pagination (omit for first page)

### 3. Frontend

```bash
cd frontend
flutter pub get
flutter run
```

The app targets `http://localhost:8080` by default. If running on a physical device or emulator, update the base URL in `lib/services/product_service.dart`:

```dart
static const String _baseUrl = 'http://<your-local-ip>:8080';
```

## How It Works

1. The Flutter app fetches a page of products from the Go backend on load.
2. Products are displayed full-screen in a vertical `PageView` (swipe up/down to navigate).
3. Each product card has a horizontal image carousel for multiple product images.
4. When you're 3 products from the end, the next page is fetched automatically using cursor-based pagination.
