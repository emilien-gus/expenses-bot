package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ai-bot-expenses/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode to suppress debug logs during testing
	gin.SetMode(gin.TestMode)
}

func TestHealthEndpoint(t *testing.T) {
	svc := service.NewExpenseService(nil)
	h := NewExpenseHandler(svc)
	router := SetupRouter(h)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expectedBody := `{"status":"ok"}`
	if w.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestCreateExpense_InvalidJSON(t *testing.T) {
	svc := service.NewExpenseService(nil)
	h := NewExpenseHandler(svc)
	router := SetupRouter(h)

	// Sending invalid/empty JSON payload
	body := bytes.NewBufferString(`{"user_id": 0}`)
	req, err := http.NewRequest(http.MethodPost, "/api/expenses", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for invalid payload, got %d", w.Code)
	}
}

func TestCORSHeaders(t *testing.T) {
	svc := service.NewExpenseService(nil)
	h := NewExpenseHandler(svc)
	router := SetupRouter(h)

	req, err := http.NewRequest(http.MethodOptions, "/api/expenses", nil)
	if err != nil {
		t.Fatalf("failed to create options request: %v", err)
	}
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check CORS preflight response
	if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
		t.Errorf("expected 200 or 204 for CORS preflight, got %d", w.Code)
	}

	allowOrigin := w.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Errorf("expected Access-Control-Allow-Origin header to be set")
	}
}
