# bridge-platform

A full-stack multiplayer contract bridge platform. Players connect via WebSocket, sit at a table, and play through the full auction and card-play phases in real time.

## Stack

- **Backend**: Go, PostgreSQL, Redis, WebSocket (gorilla/websocket)
- **Frontend**: Next.js 15, React 19, TypeScript, Tailwind CSS

## Running locally

### Prerequisites

- Go 1.22+
- Node.js
- Docker (for Postgres + Redis)

### Backend

```bash
cd backend

# Start Postgres and Redis
docker compose up -d

# Run the dev server (port 8080)
make run

# Or without Make
go run cmd/server/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### Environment variables (backend)

| Variable | Default |
|---|---|
| `PORT` | `8080` |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/bridge?sslmode=disable` |
| `REDIS_ADDR` | `localhost:6379` |

## API

`GET /health` — health check

`GET /ws` — WebSocket endpoint for game play

### WebSocket message flow

Connect to `ws://localhost:8080/ws`, then send JSON messages:

```jsonc
// 1. Create a table
{"type": "create_table"}
// → {"type": "table_created", "payload": {"tableID": "abc123"}}

// 2. Join an existing table
{"type": "join_table", "tableID": "abc123"}

// 3. Take a seat (direction: N, E, S, W)
{"type": "sit", "direction": "N"}

// 4. Start the game once all 4 seats are filled (optional seed for reproducible deals)
{"type": "start", "seed": 42}

// 5. Bid (call: P, X, XX, 1C–7NT)
{"type": "bid", "call": "1NT"}

// 6. Play a card (suit letter + rank: S=Spades, H=Hearts, D=Diamonds, C=Clubs)
{"type": "play_card", "card": "SA"}
```

After each action all seated players receive a `game_state` message with their personalized view of the board.

## Running tests

```bash
cd backend
go test ./...
```
