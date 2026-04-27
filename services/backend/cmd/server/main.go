package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"triangleoflove/backend/internal/repository"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"

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

	accountRepo := repository.NewAccountRepository(dbConn)
	authService := service.NewAuthService(accountRepo)
	healthService := service.NewHealthService(dbConn)

	checkinRepo := repository.NewCheckinRepository(dbConn)
	checkinService := service.NewCheckinService(checkinRepo)

	insightsRepo := repository.NewInsightsRepository(dbConn)
	insightsService := service.NewInsightsService(insightsRepo)

	pairingRepo := repository.NewPairingRepository(dbConn)
	coupleRepo := repository.NewCoupleRepository(dbConn)
	pairingService := service.NewPairingService(pairingRepo, coupleRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		if err := healthService.Check(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unhealthy"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	})

	mux.HandleFunc("/api/v1/status", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, statusResponse{Status: "ok", Code: http.StatusOK})
	})

	mux.HandleFunc("/api/v1/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var body struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
			return
		}

		result, err := authService.Login(r.Context(), body.Email, body.Password)
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		if err != nil {
			log.Printf("login failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}

		writeJSON(w, http.StatusOK, result)
	})

	mux.Handle("/api/v1/checkins/today", web.Middleware(web.NewCheckinHandler(checkinService)))

	mux.Handle("GET /api/v1/insights", web.Middleware(web.NewInsightsWeeklyHandler(insightsService)))
	mux.Handle("GET /api/v1/insights/{date}", web.Middleware(web.NewInsightsHandler(insightsService)))

	mux.Handle("PUT /api/v1/auth/password", web.Middleware(web.NewChangePasswordHandler(authService)))

	ph := web.NewPairingHandler(pairingService)
	mux.Handle("/api/v1/pairing", web.Middleware(http.HandlerFunc(ph.GetCode)))
	mux.Handle("/api/v1/pairing/regenerate", web.Middleware(http.HandlerFunc(ph.Regenerate)))
	mux.Handle("/api/v1/pairing/connect", web.Middleware(http.HandlerFunc(ph.Connect)))
	mux.Handle("/api/v1/couples/me", web.Middleware(http.HandlerFunc(ph.GetCoupleStatus)))
	mux.Handle("DELETE /api/v1/couples/me", web.Middleware(http.HandlerFunc(ph.Unpair)))

	mux.Handle("/api/v1/users/me", web.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		accountID, _ := r.Context().Value(web.AccountIDKey).(string)
		profile, err := authService.GetProfile(r.Context(), accountID)
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		if err != nil {
			log.Printf("get profile failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}

		writeJSON(w, http.StatusOK, profile)
	})))

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
