package service

import (
	"fmt"
	"time"

	"github.com/ai-bot-expenses/backend/internal/models"
	"github.com/ai-bot-expenses/backend/internal/repository"
)

// ExpenseService contains business logic for expense operations.
type ExpenseService struct {
	repo *repository.ExpenseRepository
}

// NewExpenseService creates a new service instance.
func NewExpenseService(repo *repository.ExpenseRepository) *ExpenseService {
	return &ExpenseService{repo: repo}
}

// Create validates input and creates a new expense.
func (s *ExpenseService) Create(req models.CreateExpenseRequest) (*models.Expense, error) {
	if !models.IsValidCategory(req.Category) {
		return nil, fmt.Errorf("invalid category: %s", req.Category)
	}

	spentAt, err := time.Parse("2006-01-02", req.SpentAt)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	expense := &models.Expense{
		UserID:   req.UserID,
		Amount:   req.Amount,
		Category: req.Category,
		SpentAt:  spentAt,
		Comment:  req.Comment,
	}

	if err := s.repo.Create(expense); err != nil {
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	return expense, nil
}

// GetByID retrieves an expense by its ID.
func (s *ExpenseService) GetByID(id uint) (*models.Expense, error) {
	return s.repo.GetByID(id)
}

// List retrieves expenses according to the provided filter.
func (s *ExpenseService) List(filter models.ExpenseFilter) ([]models.Expense, int64, error) {
	return s.repo.List(filter)
}

// Update modifies an existing expense with support for partial updates.
func (s *ExpenseService) Update(id uint, req models.UpdateExpenseRequest) (*models.Expense, error) {
	expense, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("expense not found: %w", err)
	}

	if req.UserID != nil {
		expense.UserID = *req.UserID
	}
	if req.Amount != nil {
		expense.Amount = *req.Amount
	}
	if req.Category != nil {
		if !models.IsValidCategory(*req.Category) {
			return nil, fmt.Errorf("invalid category: %s", *req.Category)
		}
		expense.Category = *req.Category
	}
	if req.SpentAt != nil {
		spentAt, err := time.Parse("2006-01-02", *req.SpentAt)
		if err != nil {
			return nil, fmt.Errorf("invalid date format: %w", err)
		}
		expense.SpentAt = spentAt
	}
	if req.Comment != nil {
		expense.Comment = *req.Comment
	}

	if err := s.repo.Update(expense); err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	return expense, nil
}

// Delete removes an expense by ID.
func (s *ExpenseService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return fmt.Errorf("expense not found: %w", err)
	}
	return s.repo.Delete(id)
}

// GetDistinctUserIDs returns all user IDs that have recorded expenses.
func (s *ExpenseService) GetDistinctUserIDs() ([]int64, error) {
	return s.repo.GetDistinctUserIDs()
}

// GetLastByUserID returns the most recent expense for a user.
func (s *ExpenseService) GetLastByUserID(userID int64) (*models.Expense, error) {
	return s.repo.GetLastByUserID(userID)
}

// FindByUserIDAndComment finds an expense by user ID and comment search term.
func (s *ExpenseService) FindByUserIDAndComment(userID int64, comment string) (*models.Expense, error) {
	return s.repo.FindByUserIDAndComment(userID, comment)
}

// SumByFilter calculates total expenses for a user with the given filter.
func (s *ExpenseService) SumByFilter(userID int64, filter models.ExpenseFilter) (float64, error) {
	return s.repo.SumByFilter(userID, filter)
}
