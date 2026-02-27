# Texas Holdem API
A simple Texas Holdem Poker Game played by 2 players



## Architecture

```
texas_holdem_api
|-- go.mod
|-- handlers.go     Handler functions that run on each request
├── main.go         Main Server application
|-- middleware.go   Configures CORS
└── types.go        Includes Types
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.21+

### Run the Backend

```bash
go run .
```

The server starts on `http://localhost:8080`.

### Build
```bash
go build
```

### Run the build
```bash
./texas_holdem_api
```

## API

| Protocol | Method | Endpoint               | Description       |
| -------- | ------ | ---------------------- |
| HTTP     | GET    | /api/play/texas-holdem | Gets player cards |

### GET /api/play/texas-holdem

Sample Response
```json
{ 
    "player1": ["SK", "CA"],
    "player2": ["HA", "SQ"],
    "communityCards": ["D6", "S9", "H4" , "S3", "C2"]
}
```
