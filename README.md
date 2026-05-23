# Ticket Tracker API
A simple REST API for tracking tickets, built using Go.

## How to Run
Make sure you have Go installed, then:
```
go run main.go
```
The server starts on port 8080.

## Endpoints
**Create a ticket**
```
POST /tickets
```
Body:
```json
{
  "title": "Fix login bug",
  "description": "Users cant log in on mobile"
}
```

**Get all tickets**
```
GET /tickets
```

**Get one ticket**
```
GET /tickets/1
```

## Design Decisions
- Mutex used to prevent race conditions across concurrent requests
- Status and ID are set server-side to prevent client manipulation
- In-memory storage intentionally kept simple — PostgreSQL version in progress
