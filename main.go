// main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/go-chi/chi/v5"
)

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

var (
	products []Product
	index    bleve.Index
)

func generateProducts(n int) []Product {
	categories := []string{"Electronics", "Footwear", "Fitness", "Home & Kitchen", "Books", "Toys"}
	var items []Product
	for i := 1; i <= n; i++ {
		item := Product{
			ID:       i,
			Name:     fmt.Sprintf("Product %d", i),
			Category: categories[rand.Intn(len(categories))],
		}
		items = append(items, item)
	}
	return items
}

func indexProducts(products []Product) bleve.Index {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.NewMemOnly(mapping)
	if err != nil {
		log.Fatalf("Error creating Bleve index: %v", err)
	}
	for _, p := range products {
		err := idx.Index(strconv.Itoa(p.ID), p)
		if err != nil {
			log.Printf("Failed to index product ID %d: %v", p.ID, err)
		}
	}
	return idx
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query param 'q' is required", http.StatusBadRequest)
		return
	}
	q := bleve.NewQueryStringQuery(query)
	sreq := bleve.NewSearchRequestOptions(q, 50, 0, false)
	res, err := index.Search(sreq)
	if err != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}
	var results []Product
	for _, hit := range res.Hits {
		id, err := strconv.Atoi(hit.ID)
		if err != nil || id <= 0 || id > len(products) {
			continue
		}
		results = append(results, products[id-1])
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	products = generateProducts(1_000_000)
	index = indexProducts(products)

	r := chi.NewRouter()
	r.Get("/search", searchHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Println("Server listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly")
}
