# bridge-platform

## Backend

Requires [Go 1.22+](https://go.dev/dl/).

```bash
cd backend

# Run the server (development)
make run

# Or without Make:
go run cmd/server/main.go

# Build a compiled binary
make build        # outputs to backend/bin/server

# Run the compiled binary
./bin/server
```

## Frontend

```
cd frontend
npm install
npm run dev
```