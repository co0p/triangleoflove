package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"triangleoflove/backend/internal/repository"
	"triangleoflove/backend/internal/service"

	_ "github.com/lib/pq"
)

type statusResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func main() {
	dsn, err := postgresDSN()
	if err != nil {
		log.Fatalf("database configuration error: %v", err)
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer dbConn.Close()

	if err := waitForDatabase(dbConn, 30, 2*time.Second); err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	repo := repository.NewRoundtripRepository(dbConn)
	roundtripService := service.NewRoundtripService(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, statusResponse{Status: "ok", Code: http.StatusOK})
	})

	mux.HandleFunc("/demo/roundtrip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		result, err := roundtripService.Execute(r.Context())
		if err != nil {
			log.Printf("roundtrip execution failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}

		writeJSON(w, http.StatusOK, result)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func postgresDSN() (string, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return "", errMissingEnv("DATABASE_URL")
	}

	return dsn, nil
}

func errMissingEnv(key string) error {
	return &missingEnvError{Key: key}
}

type missingEnvError struct {
	Key string
}

func (e *missingEnvError) Error() string {
	return e.Key + " is required"
}

func waitForDatabase(dbConn *sql.DB, maxAttempts int, delay time.Duration) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := dbConn.Ping(); err == nil {
			return nil
		} else {
			lastErr = err
			log.Printf("db not ready (attempt %d/%d): %v", attempt, maxAttempts, err)
		}

		time.Sleep(delay)
	}

	return lastErr
}
