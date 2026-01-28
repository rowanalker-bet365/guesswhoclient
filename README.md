# Guess Who Client API

This is the client-facing Golang API for the Guess Who UI application.

## Project Structure

```
guesswhoclientapi/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handlers.go
│   ├── data/
│   │   ├── models.go
│   │   └── store.go
│   ├── redis/
│   │   └── pubsub.go
│   └── sse/
│       └── sse.go
├── go.mod
└── go.sum
```

## How to Run the API Server

1.  Navigate to the `guesswhoclientapi` directory:
    ```bash
    cd guesswhoclientapi
    ```
2.  Run the main application:
    ```bash
    go run guesswhoclientapi.go
    ```
3.  The API server will start on `http://localhost:8080`.

## Redis Integration (Optional)

This API can connect to a Redis server to listen for real-time game updates via a Pub/Sub channel. This allows the API to be scaled horizontally while still providing real-time updates to all connected clients.

To enable this feature, set the `REDIS_ADDR` environment variable to the address of your Redis server (e.g., `localhost:6379`).

```bash
export REDIS_ADDR=localhost:6379
go run guesswhoclientapi.go
```

If `REDIS_ADDR` is not set, the application will start normally and use only its in-memory data store, without attempting to connect to Redis. When connected, the API subscribes to the `game_updates` channel. Any message published to this channel will trigger a broadcast to all connected SSE clients, instructing them to refetch the latest game state.

## API Endpoints

*   `POST /api/auth/signup`: Creates a new team.
*   `POST /api/auth/login`: Authenticates a team and returns a JWT.
*   `GET /api/game/state`: Returns the public game state.
*   `GET /api/team/progress`: Returns the private progress for the authenticated team.
*   `POST /api/team/reset`: Resets the `solvedCharacters` array for the authenticated team.
*   `POST /api/teams/update-data`: Simulates a game data update and broadcasts an SSE event.
*   `GET /events`: The Server-Sent Events endpoint for real-time updates.

## SSE Functionality

You can test the SSE functionality by navigating to `http://localhost:8080/events` in your browser and then hitting the `http://localhost:8080/api/teams/update-data` endpoint (e.g. with `curl` or Postman) to see the update message being broadcast.