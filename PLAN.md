# Product Feeds App - Implementation Plan

## Overview
A TikTok-style vertical swipe feed for browsing products. Flutter frontend, Go backend, Supabase for storage.

## Current State
- **Supabase** is set up with a `products` table (28 rows) and a `product-images` public storage bucket
- **Product schema**: `id` (uuid), `name`, `price`, `currency`, `image_paths` (jsonb array of URLs), `source_url`, `company_name`, `scraped_at`
- Each product has 2+ images stored as public Supabase storage URLs
- No backend or frontend code exists yet

---

## Phase 1: Go Backend API

### 1.1 Project Setup
- Initialize Go module at `backend/`
- Dependencies: `github.com/gin-gonic/gin` (HTTP router), `github.com/joho/godotenv` (.env loading), `github.com/jackc/pgx/v5` (Postgres driver)
- Load config from `.env` (Supabase DB connection string derived from project URL + password)

### 1.2 Database Connection
- Connect to Supabase Postgres directly: `postgresql://postgres.[project-ref]:[password]@aws-0-eu-west-2.pooler.supabase.com:6543/postgres`
- Create a `db` package with connection pool setup

### 1.3 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/products` | Paginated product feed (cursor-based) |
| `GET` | `/api/products/:id` | Single product detail |

#### `GET /api/products`
- **Query params**: `cursor` (uuid, optional), `limit` (int, default 10)
- **Cursor-based pagination** using `id` — returns products where `id > cursor`, ordered by `scraped_at DESC`
- **Response**:
  ```json
  {
    "products": [...],
    "next_cursor": "uuid-of-last-item"
  }
  ```

### 1.4 Product Model
```go
type Product struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Price       float64  `json:"price"`
    Currency    string   `json:"currency"`
    ImagePaths  []string `json:"image_paths"`
    SourceURL   string   `json:"source_url"`
    CompanyName string   `json:"company_name"`
    ScrapedAt   string   `json:"scraped_at"`
}
```

### 1.5 File Structure
```
backend/
├── main.go              # Entry point, router setup
├── go.mod
├── .env -> ../.env      # Symlink to root .env
├── config/
│   └── config.go        # Env loading
├── db/
│   └── db.go            # Postgres connection pool
├── handlers/
│   └── products.go      # HTTP handlers
└── models/
    └── product.go       # Product struct
```

---

## Phase 2: Flutter Frontend

### 2.1 Project Setup
- Create Flutter project at `frontend/`
- Dependencies: `http`, `cached_network_image`, `flutter_riverpod` (state management)

### 2.2 Core UX: Vertical Swipe Feed
- Use `PageView` with `scrollDirection: Axis.vertical` — this gives the TikTok-style full-screen snap-to-card swiping
- Each "page" is one product, filling the entire screen

### 2.3 Product Card (Full-Screen Page)
Each card shows:
- **Product image** — full-screen background, covers the card. If multiple images, horizontal swipe/carousel with dots indicator
- **Bottom overlay** (gradient fade to dark):
  - Product name
  - Price + currency
  - Company name
  - "View Source" button (opens `source_url` in browser)

### 2.4 Image Carousel
- Within each product card, use a horizontal `PageView` for the `image_paths` array
- Dot indicators at bottom showing current image position
- Preload/cache adjacent images with `cached_network_image`

### 2.5 Infinite Scroll / Pagination
- Load initial batch of 10 products
- When user is 3 products from the end of loaded list, trigger next page fetch using cursor
- Show a loading shimmer while fetching

### 2.6 State Management
- Use Riverpod with a `ProductFeedNotifier` (AsyncNotifier)
- State holds: `List<Product>`, `nextCursor`, `isLoading`, `hasMore`
- Provider triggers pagination automatically

### 2.7 File Structure
```
frontend/
├── lib/
│   ├── main.dart
│   ├── config.dart              # API base URL
│   ├── models/
│   │   └── product.dart         # Product model + fromJson
│   ├── services/
│   │   └── product_service.dart # HTTP client for API
│   ├── providers/
│   │   └── feed_provider.dart   # Riverpod feed state
│   └── screens/
│       └── feed_screen.dart     # Main vertical feed page
│       └── widgets/
│           ├── product_card.dart     # Full-screen product card
│           └── image_carousel.dart   # Horizontal image swiper
├── pubspec.yaml
└── ...
```

---

## Phase 3: Polish & Enhancements (Optional / Later)

- **Like/Save**: Add a `user_likes` table, heart button on cards, saved products screen
- **Category filtering**: Add categories to products, filter chips at top
- **Share**: Native share sheet for product links
- **Animations**: Smooth transitions between cards, parallax on images
- **Search**: Search bar with product name / company autocomplete
- **Auth**: Supabase Auth for user accounts

---

## Build & Run Order

1. **Backend first** — get the API serving products
2. **Flutter feed screen** — wire up to API, get vertical swiping working
3. **Image carousel** — add horizontal image swiping within each card
4. **Pagination** — infinite scroll with cursor
5. **Visual polish** — gradients, typography, loading states
