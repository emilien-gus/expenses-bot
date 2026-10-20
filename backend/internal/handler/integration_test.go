package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ai-bot-expenses/backend/internal/handler"
	"github.com/ai-bot-expenses/backend/internal/models"
	"github.com/ai-bot-expenses/backend/internal/repository"
	"github.com/ai-bot-expenses/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type singleResponse struct {
	Data models.Expense `json:"data"`
}

type listResponse struct {
	Data  []models.Expense `json:"data"`
	Total int64            `json:"total"`
}

// setupIntegrationTest initializes an in-memory SQLite database and full HTTP router.
func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	dbName := fmt.Sprintf("file:mem_%p?mode=memory&cache=shared", t)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory sqlite: %v", err)
	}

	if err := db.AutoMigrate(&models.Expense{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	repo := repository.NewExpenseRepository(db)
	svc := service.NewExpenseService(repo)
	h := handler.NewExpenseHandler(svc)
	router := handler.SetupRouter(h)

	return router, db
}

func TestIntegration_ExpenseLifecycleCRUD(t *testing.T) {
	router, _ := setupIntegrationTest(t)

	// 1. CREATE Expense via POST /api/expenses
	createReq := models.CreateExpenseRequest{
		UserID:   777,
		Amount:   350.50,
		Category: "еда",
		SpentAt:  "2026-10-19",
		Comment:  "Кофе с круассаном",
	}
	body, _ := json.Marshal(createReq)

	req, _ := http.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("POST /api/expenses expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var createResp singleResponse
	if err := json.Unmarshal(w.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("failed to unmarshal created expense: %v", err)
	}
	created := createResp.Data
	if created.ID == 0 {
		t.Fatal("expected non-zero ID for created expense")
	}
	if created.Amount != 350.50 || created.Category != "еда" || created.UserID != 777 {
		t.Errorf("unexpected expense fields: %+v", created)
	}

	// 2. READ list via GET /api/expenses?user_id=777
	req, _ = http.NewRequest(http.MethodGet, "/api/expenses?user_id=777", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/expenses expected 200 OK, got %d", w.Code)
	}

	var listResp listResponse
	if err := json.Unmarshal(w.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("failed to unmarshal list response: %v", err)
	}
	if listResp.Total != 1 || len(listResp.Data) != 1 {
		t.Fatalf("expected 1 item, got total=%d len=%d", listResp.Total, len(listResp.Data))
	}

	// 3. READ by ID via GET /api/expenses/:id
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/expenses/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/expenses/%d expected 200, got %d", created.ID, w.Code)
	}

	// 4. UPDATE via PUT /api/expenses/:id
	newAmount := 450.00
	newCategory := "покупки"
	newSpentAt := "2026-10-19"
	newComment := "Кофе, круассан и вода"
	updateReq := models.UpdateExpenseRequest{
		Amount:   &newAmount,
		Category: &newCategory,
		SpentAt:  &newSpentAt,
		Comment:  &newComment,
	}
	updateBody, _ := json.Marshal(updateReq)
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/expenses/%d", created.ID), bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("PUT /api/expenses/%d expected 200 OK, got %d: %s", created.ID, w.Code, w.Body.String())
	}

	var updateResp singleResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updateResp); err != nil {
		t.Fatalf("failed to unmarshal updated response: %v", err)
	}
	updated := updateResp.Data
	if updated.Amount != 450.00 || updated.Category != "покупки" {
		t.Errorf("updated fields mismatch: %+v", updated)
	}

	// 5. GET USERS via GET /api/users
	req, _ = http.NewRequest(http.MethodGet, "/api/users", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET /api/users expected 200, got %d", w.Code)
	}
	var usersResp struct {
		Data []int64 `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &usersResp); err != nil {
		t.Fatalf("failed to unmarshal users response: %v", err)
	}
	if len(usersResp.Data) != 1 || usersResp.Data[0] != 777 {
		t.Errorf("expected users [777], got %+v", usersResp)
	}

	// 6. DELETE via DELETE /api/expenses/:id
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/expenses/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("DELETE /api/expenses/%d expected 200 OK, got %d", created.ID, w.Code)
	}

	// 7. VERIFY DELETED via GET /api/expenses/:id -> 404
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/expenses/%d", created.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found after deletion, got %d", w.Code)
	}
}

func TestIntegration_ExpenseFiltering(t *testing.T) {
	router, _ := setupIntegrationTest(t)

	// Create 3 expenses
	expenses := []models.CreateExpenseRequest{
		{UserID: 100, Amount: 100, Category: "еда", SpentAt: "2026-10-10", Comment: "Обед"},
		{UserID: 100, Amount: 200, Category: "транспорт", SpentAt: "2026-10-15", Comment: "Метро"},
		{UserID: 200, Amount: 500, Category: "еда", SpentAt: "2026-10-18", Comment: "Ужин"},
	}

	for _, exp := range expenses {
		body, _ := json.Marshal(exp)
		req, _ := http.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create test expense: %s", w.Body.String())
		}
	}

	// Filter by UserID: 100
	req, _ := http.NewRequest(http.MethodGet, "/api/expenses?user_id=100", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp listResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 2 {
		t.Errorf("expected 2 items for user 100, got %d", resp.Total)
	}

	// Filter by Category: "еда"
	req, _ = http.NewRequest(http.MethodGet, "/api/expenses?category=еда", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 2 {
		t.Errorf("expected 2 items for category 'еда', got %d", resp.Total)
	}

	// Filter by UserID: 100 and Category: "еда"
	req, _ = http.NewRequest(http.MethodGet, "/api/expenses?user_id=100&category=еда", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 1 {
		t.Errorf("expected 1 item for user 100 with category 'еда', got %d", resp.Total)
	}

	// Filter by Date Range: 2026-10-14 to 2026-10-19
	req, _ = http.NewRequest(http.MethodGet, "/api/expenses?date_from=2026-10-14&date_to=2026-10-19", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 2 {
		t.Errorf("expected 2 items in date range, got %d", resp.Total)
	}
}
