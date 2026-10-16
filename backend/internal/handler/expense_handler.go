package handler

import (
	"net/http"
	"strconv"

	"github.com/ai-bot-expenses/backend/internal/models"
	"github.com/ai-bot-expenses/backend/internal/service"
	"github.com/gin-gonic/gin"
)

// ExpenseHandler handles HTTP requests for expense operations.
type ExpenseHandler struct {
	service *service.ExpenseService
}

// NewExpenseHandler creates a new handler instance.
func NewExpenseHandler(svc *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: svc}
}

// RegisterRoutes sets up all expense-related routes on the given router group.
func (h *ExpenseHandler) RegisterRoutes(rg *gin.RouterGroup) {
	expenses := rg.Group("/expenses")
	{
		expenses.GET("", h.List)
		expenses.POST("", h.Create)
		expenses.GET("/:id", h.GetByID)
		expenses.PUT("/:id", h.Update)
		expenses.DELETE("/:id", h.Delete)
	}

	rg.GET("/users", h.ListUsers)
}

// List returns expenses with optional filters.
// GET /api/expenses?user_id=123&category=еда&date_from=2026-10-01&date_to=2026-10-31
func (h *ExpenseHandler) List(c *gin.Context) {
	var filter models.ExpenseFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters: " + err.Error()})
		return
	}

	expenses, total, err := h.service.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list expenses: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  expenses,
		"total": total,
	})
}

// Create adds a new expense.
// POST /api/expenses
func (h *ExpenseHandler) Create(c *gin.Context) {
	var req models.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	expense, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": expense})
}

// GetByID returns a single expense by ID.
// GET /api/expenses/:id
func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	expense, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense})
}

// Update modifies an existing expense.
// PUT /api/expenses/:id
func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	var req models.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	expense, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expense})
}

// Delete removes an expense by ID.
// DELETE /api/expenses/:id
func (h *ExpenseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense deleted"})
}

// ListUsers returns all unique user IDs for filtering.
// GET /api/users
func (h *ExpenseHandler) ListUsers(c *gin.Context) {
	userIDs, err := h.service.GetDistinctUserIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userIDs})
}
