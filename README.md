# 🏍️ SpotSync — Parking Zone API

A production-ready REST API for managing parking zones and reservations, built with **Go**, **Echo v5**, and **GORM**.

## live link : [https://spotsync-pqf9.onrender.com/]

---

## 📦 Tech Stack

| Layer     | Technology         |
| --------- | ------------------ |
| Language  | Go 1.22+           |
| Framework | Echo v5            |
| ORM       | GORM               |
| Database  | PostgreSQL         |
| Auth      | JWT (`golang-jwt`) |
| Password  | bcrypt             |

---

## 🗂️ Project Structure

```
SpotSync/
├── cmd/
│   └── main.go                  # Entry point
├── internal/
│   ├── auth/
│   │   └── jwt.go               # JWT service
│   ├── config/
│   │   └── config.go            # Env config loader
│   ├── domain/
│   │   ├── users/               # User domain
│   │   ├── parking_zones/       # Parking zone domain
│   │   │   └── dto/
│   │   └── reservations/        # Reservations domain
│   │       └── dto/
│   ├── http_response/
│   │   └── error.go             # Shared error response struct
│   └── middlewares/
│       ├── auth.go              # JWT auth middleware
│       └── admin.go             # Admin-only middleware
```

---

## ⚙️ Environment Variables

Create a `.env` file in the project root:

```env
DNS=postgres://username:password@localhost:5432/spotsync?sslmode=disable
PORT=8080
JWT_SECRET_KEY=your-super-secret-key
```

---

## 🚀 Getting Started

```bash
# 1. Clone the repository
git clone https://github.com/asraful-07/SpotSync
cd spotsync

# 2. Install dependencies
go mod tidy

# 3. Set up environment variables
cp .env.example .env

# 4. Run the server
go run cmd/main.go
```

The server starts at `http://localhost:8080`.

---

## 🔐 Authentication

All protected routes require a Bearer token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

Roles: `driver` · `admin`

---

## 📡 API Reference

### Auth

| Method | Endpoint                | Access | Description           |
| ------ | ----------------------- | ------ | --------------------- |
| `POST` | `/api/v1/auth/register` | Public | Register a new user   |
| `POST` | `/api/v1/auth/login`    | Public | Login and receive JWT |

---

### 🅿️ Parking Zones

| Method | Endpoint            | Access | Description                           |
| ------ | ------------------- | ------ | ------------------------------------- |
| `POST` | `/api/v1/zones`     | Admin  | Create a parking zone                 |
| `GET`  | `/api/v1/zones`     | Public | List all zones with live availability |
| `GET`  | `/api/v1/zones/:id` | Public | Get a single zone                     |

**Create Zone — Request**

```json
{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 20,
  "price_per_hour": 5.5
}
```

> `type` must be one of: `general` · `ev_charging` · `covered`

**Get All Zones — Response**

```json
{
  "success": true,
  "message": "Parking zones retrieved successfully",
  "data": [
    {
      "id": 5,
      "name": "Terminal 1 EV Charging",
      "type": "ev_charging",
      "total_capacity": 20,
      "available_spots": 14,
      "price_per_hour": 5.5,
      "created_at": "2026-06-20T10:30:00Z"
    }
  ]
}
```

> `available_spots` is calculated dynamically: `total_capacity − count(active reservations)`.

---

### 🎟️ Reservations

| Method   | Endpoint                               | Access         | Description            |
| -------- | -------------------------------------- | -------------- | ---------------------- |
| `POST`   | `/api/v1/reservations`                 | Driver · Admin | Reserve a parking spot |
| `GET`    | `/api/v1/reservations/my-reservations` | Driver · Admin | Get own reservations   |
| `DELETE` | `/api/v1/reservations/:id`             | Driver · Admin | Cancel a reservation   |
| `GET`    | `/api/v1/reservations`                 | Admin          | Get all reservations   |

**Reserve a Spot — Request**

```json
{
  "zone_id": 5,
  "license_plate": "ABC-1234"
}
```

**Reserve a Spot — Response `201`**

```json
{
  "success": true,
  "message": "Reservation confirmed successfully",
  "data": {
    "id": 105,
    "user_id": 1,
    "zone_id": 8,
    "license_plate": "ABC-1234",
    "status": "active",
    "created_at": "2026-06-20T15:30:00Z",
    "updated_at": "2026-06-20T15:30:00Z"
  }
}
```

**Get My Reservations — Response `200`**

```json
{
  "success": true,
  "message": "My reservations retrieved successfully",
  "data": [
    {
      "id": 10,
      "license_plate": "ABC-1234",
      "status": "active",
      "zone": {
        "id": 8,
        "name": "Terminal 1 EV Charging",
        "type": "ev_charging"
      },
      "created_at": "2026-06-20T15:30:00Z"
    }
  ]
}
```

**Cancel Reservation — Response `200`**

```json
{
  "success": true,
  "message": "Reservation cancelled successfully"
}
```

---

## ⚠️ Concurrency & Capacity Safety

The reservation endpoint is protected against race conditions using **database-level row locking**.

When two drivers attempt to book the last available spot simultaneously:

```
Request A                        Request B
─────────────────────────────────────────────────────
BEGIN TRANSACTION                BEGIN TRANSACTION
SELECT zone FOR UPDATE  ✅ lock  SELECT zone FOR UPDATE  ⏳ waits
COUNT active = 19
19 < 20 capacity ✅
INSERT reservation
COMMIT  →  releases lock         ✅ gets lock
                                 COUNT active = 20
                                 20 >= 20 capacity ❌
                                 → 409 Zone Full
```

Implementation uses `gorm.io/gorm/clause.Locking{Strength: "UPDATE"}` inside a `db.Transaction(...)`.

---

## 🛑 Error Responses

All errors follow a consistent shape:

```json
{
  "code": 404,
  "message": "Parking zone not found"
}
```

| HTTP Code | Meaning                                 |
| --------- | --------------------------------------- |
| `400`     | Bad request / validation failed         |
| `401`     | Missing or invalid JWT                  |
| `403`     | Forbidden — you don't own this resource |
| `404`     | Resource not found                      |
| `409`     | Zone at full capacity                   |
| `500`     | Internal server error                   |

---

## 🗄️ Database Models

### `parking_zones`

| Column           | Type         | Notes                               |
| ---------------- | ------------ | ----------------------------------- |
| `id`             | uint         | Primary key                         |
| `name`           | varchar(255) | Required                            |
| `type`           | varchar(100) | `general`, `ev_charging`, `covered` |
| `total_capacity` | int          | Required, > 0                       |
| `price_per_hour` | float64      | Required, > 0                       |
| `created_at`     | timestamp    | Auto                                |
| `updated_at`     | timestamp    | Auto                                |
| `deleted_at`     | timestamp    | Soft delete                         |

### `reservations`

| Column          | Type        | Notes                              |
| --------------- | ----------- | ---------------------------------- |
| `id`            | uint        | Primary key                        |
| `user_id`       | uint        | FK → users                         |
| `zone_id`       | uint        | FK → parking_zones                 |
| `license_plate` | varchar(15) | Required                           |
| `status`        | varchar(20) | `active`, `completed`, `cancelled` |
| `created_at`    | timestamp   | Auto                               |
| `updated_at`    | timestamp   | Auto                               |
| `deleted_at`    | timestamp   | Soft delete                        |

---

## 📝 License

MIT © 2026 SpotSync
