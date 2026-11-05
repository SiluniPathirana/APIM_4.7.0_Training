// ABC Financial mock back ends for Lab 04 (Creating an API Product).
//
// One program starts two separate back ends, one per lab API:
//
//	Customer Info back end - http://localhost:8091 (customer-info-api.yaml)
//	  GET    /customers
//	  POST   /customers
//	  GET    /customers/{customerId}
//	  PUT    /customers/{customerId}
//	  DELETE /customers/{customerId}
//
//	Leasing back end - http://localhost:8092 (leasing-api.yaml)
//	  GET    /assets
//	  POST   /assets
//	  GET    /assets/{assetId}
//	  PUT    /assets/{assetId}
//	  DELETE /assets/{assetId}
//
// Data is kept in memory and resets when the program restarts.
//
// Run:   go run main.go            (or one of the pre-built binaries)
// Ports: go run main.go -customer-port 9091 -leasing-port 9092
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Customer struct {
	ID          string  `json:"customerId"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	AccountType string  `json:"accountType"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency"`
}

type Asset struct {
	ID            string  `json:"assetId"`
	CustomerID    string  `json:"customerId"`
	Type          string  `json:"type"`
	Description   string  `json:"description"`
	LeaseAmount   float64 `json:"leaseAmount"`
	MonthlyRental float64 `json:"monthlyRental"`
	TermMonths    int     `json:"termMonths"`
	StartDate     string  `json:"startDate"`
	Status        string  `json:"status"`
}

type store struct {
	mu        sync.Mutex
	customers map[string]Customer
	assets    map[string]Asset
	nextCust  int
	nextAsset int
}

func newStore() *store {
	s := &store{
		customers: map[string]Customer{},
		assets:    map[string]Asset{},
		nextCust:  1004,
		nextAsset: 5004,
	}
	for _, c := range []Customer{
		{"1001", "Jane Wanjiku", "jane.wanjiku@example.com", "+254700000001", "SAVINGS", 152340.50, "KES"},
		{"1002", "Peter Otieno", "peter.otieno@example.com", "+254700000002", "CURRENT", 98400.00, "KES"},
		{"1003", "Amina Hassan", "amina.hassan@example.com", "+254700000003", "BUSINESS", 1250000.00, "KES"},
	} {
		s.customers[c.ID] = c
	}
	for _, a := range []Asset{
		{"5001", "1001", "VEHICLE", "Toyota Axio 2021", 1800000, 62500, 36, "2025-02-01", "ACTIVE"},
		{"5002", "1003", "EQUIPMENT", "Industrial printer", 950000, 44000, 24, "2025-06-15", "ACTIVE"},
		{"5003", "1002", "VEHICLE", "Isuzu D-Max 2022", 3200000, 98000, 48, "2024-11-01", "CLOSED"},
	} {
		s.assets[a.ID] = a
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func notFound(w http.ResponseWriter, kind, id string) {
	writeJSON(w, http.StatusNotFound, map[string]string{
		"code": "404", "message": fmt.Sprintf("%s %s not found", kind, id),
	})
}

func badRequest(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{
		"code": "400", "message": "invalid JSON body: " + err.Error(),
	})
}

// ---- Customer Info API ----

func (s *store) listCustomers(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := make([]Customer, 0, len(s.customers))
	for _, c := range s.customers {
		list = append(list, c)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	writeJSON(w, http.StatusOK, map[string]any{"count": len(list), "customers": list})
}

func (s *store) createCustomer(w http.ResponseWriter, r *http.Request) {
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		badRequest(w, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = strconv.Itoa(s.nextCust)
	s.nextCust++
	if c.Currency == "" {
		c.Currency = "KES"
	}
	s.customers[c.ID] = c
	writeJSON(w, http.StatusCreated, c)
}

func (s *store) getCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("customerId")
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.customers[id]
	if !ok {
		notFound(w, "customer", id)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *store) updateCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("customerId")
	var c Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		badRequest(w, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.customers[id]; !ok {
		notFound(w, "customer", id)
		return
	}
	c.ID = id
	s.customers[id] = c
	writeJSON(w, http.StatusOK, c)
}

func (s *store) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("customerId")
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.customers[id]; !ok {
		notFound(w, "customer", id)
		return
	}
	delete(s.customers, id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "customer " + id + " deleted"})
}

// ---- Leasing API ----

func (s *store) listAssets(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := make([]Asset, 0, len(s.assets))
	for _, a := range s.assets {
		list = append(list, a)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	writeJSON(w, http.StatusOK, map[string]any{"count": len(list), "assets": list})
}

func (s *store) createAsset(w http.ResponseWriter, r *http.Request) {
	var a Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		badRequest(w, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	a.ID = strconv.Itoa(s.nextAsset)
	s.nextAsset++
	if a.Status == "" {
		a.Status = "PENDING"
	}
	if a.StartDate == "" {
		a.StartDate = time.Now().Format("2006-01-02")
	}
	s.assets[a.ID] = a
	writeJSON(w, http.StatusCreated, a)
}

func (s *store) getAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("assetId")
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.assets[id]
	if !ok {
		notFound(w, "asset", id)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (s *store) updateAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("assetId")
	var a Asset
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		badRequest(w, err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.assets[id]; !ok {
		notFound(w, "asset", id)
		return
	}
	a.ID = id
	s.assets[id] = a
	writeJSON(w, http.StatusOK, a)
}

func (s *store) deleteAsset(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("assetId")
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.assets[id]; !ok {
		notFound(w, "asset", id)
		return
	}
	delete(s.assets, id)
	writeJSON(w, http.StatusOK, map[string]string{"message": "asset " + id + " deleted"})
}

// logRequests prints one line per call, tagged with the back end name, so
// attendees can see which back end the gateway sends each request to.
func logRequests(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s (%s)", name, r.Method, r.URL.Path, time.Since(start).Round(time.Microsecond))
	})
}

func health(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"service": name, "status": "UP"})
	}
}

func serve(name string, port int, mux *http.ServeMux, errs chan<- error) {
	addr := fmt.Sprintf(":%d", port)
	log.Printf("%s back end listening on http://localhost%s", name, addr)
	errs <- fmt.Errorf("%s back end: %w", name, http.ListenAndServe(addr, logRequests(name, mux)))
}

func main() {
	customerPort := flag.Int("customer-port", 8091, "port of the Customer Info back end")
	leasingPort := flag.Int("leasing-port", 8092, "port of the Leasing back end")
	flag.Parse()

	s := newStore()

	customers := http.NewServeMux()
	customers.HandleFunc("GET /customers", s.listCustomers)
	customers.HandleFunc("POST /customers", s.createCustomer)
	customers.HandleFunc("GET /customers/{customerId}", s.getCustomer)
	customers.HandleFunc("PUT /customers/{customerId}", s.updateCustomer)
	customers.HandleFunc("DELETE /customers/{customerId}", s.deleteCustomer)
	customers.HandleFunc("GET /{$}", health("Customer Info"))

	leasing := http.NewServeMux()
	leasing.HandleFunc("GET /assets", s.listAssets)
	leasing.HandleFunc("POST /assets", s.createAsset)
	leasing.HandleFunc("GET /assets/{assetId}", s.getAsset)
	leasing.HandleFunc("PUT /assets/{assetId}", s.updateAsset)
	leasing.HandleFunc("DELETE /assets/{assetId}", s.deleteAsset)
	leasing.HandleFunc("GET /{$}", health("Leasing"))

	// If either back end cannot start (for example its port is in use), stop the program.
	errs := make(chan error, 2)
	go serve("Customer Info", *customerPort, customers, errs)
	go serve("Leasing", *leasingPort, leasing, errs)
	log.Fatal(<-errs)
}
