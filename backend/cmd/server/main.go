package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	dbConn, err := sql.Open("postgres", postgresDSN())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer dbConn.Close()

	if err := dbConn.Ping(); err != nil {
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

func postgresDSN() string {
	host := envOrDefault("POSTGRES_HOST", "db")
	port := envOrDefault("POSTGRES_PORT", "5432")
	user := envOrDefault("POSTGRES_USER", "triangle")
	password := envOrDefault("POSTGRES_PASSWORD", "triangle")
	dbName := envOrDefault("POSTGRES_DB", "triangleoflove")

	return "host=" + host +
		" port=" + port +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbName +
		" sslmode=disable"
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
