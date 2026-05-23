package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Ticket struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type Store struct {
	mu      sync.Mutex
	tickets []Ticket
	nextID  int
}

func NewStore() *Store {
	return &Store{nextID: 1, tickets: []Ticket{}}
}

func (s *Store) createTicket(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ticket.ID = s.nextID
	ticket.Status = "open"
	ticket.CreatedAt = time.Now()
	s.nextID++
	s.tickets = append(s.tickets, ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ticket)
}

func (s *Store) getTickets(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.tickets)
}

func (s *Store) getTicket(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tickets/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, ticket := range s.tickets {
		if ticket.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ticket)
			return
		}
	}

	http.Error(w, "Ticket not found", http.StatusNotFound)
}

func main() {
	store := NewStore()

	http.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			store.getTickets(w, r)
		case http.MethodPost:
			store.createTicket(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tickets/", func(w http.ResponseWriter, r *http.Request) {
		store.getTicket(w, r)
	})

	http.ListenAndServe(":8080", nil)
}
