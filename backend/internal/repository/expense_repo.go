package repository

import (
	"time"

	"github.com/ai-bot-expenses/backend/internal/models"
	"gorm.io/gorm"
)

// ExpenseRepository handles all database operations for expenses.
type ExpenseRepository struct {
	db *gorm.DB
}

// NewExpenseRepository creates a new repository instance.
func NewExpenseRepository(db *gorm.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create inserts a new expense into the database.
func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

// GetByID retrieves a single expense by its ID.
func (r *ExpenseRepository) GetByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.First(&expense, id).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// List retrieves expenses with optional filters and total count.
func (r *ExpenseRepository) List(filter models.ExpenseFilter) ([]models.Expense, int64, error) {
	query := r.db.Model(&models.Expense{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", *filter.DateFrom)
		if err == nil {
			query = query.Where("spent_at >= ?", dateFrom)
		}
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", *filter.DateTo)
		if err == nil {
			query = query.Where("spent_at <= ?", dateTo)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var expenses []models.Expense
	err := query.Order("spent_at DESC, id DESC").Find(&expenses).Error
	return expenses, total, err
}

// ListByUserID retrieves all expenses for a specific user.
func (r *ExpenseRepository) ListByUserID(userID int64) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.Where("user_id = ?", userID).Order("spent_at DESC, id DESC").Find(&expenses).Error
	return expenses, err
}

// GetLastByUserID retrieves the most recent expense created by a user.
func (r *ExpenseRepository) GetLastByUserID(userID int64) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// FindByUserIDAndComment searches for an expense by user ID and comment substring.
func (r *ExpenseRepository) FindByUserIDAndComment(userID int64, comment string) (*models.Expense, error) {
	var expense models.Expense
	err := r.db.Where("user_id = ? AND LOWER(comment) LIKE LOWER(?)", userID, "%"+comment+"%").
		Order("created_at DESC").First(&expense).Error
	if err != nil {
		return nil, err
	}
	return &expense, nil
}

// Update saves changes to an existing expense.
func (r *ExpenseRepository) Update(expense *models.Expense) error {
	return r.db.Save(expense).Error
}

// Delete removes an expense by its ID.
func (r *ExpenseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Expense{}, id).Error
}

// GetDistinctUserIDs returns all unique user IDs that have expenses.
func (r *ExpenseRepository) GetDistinctUserIDs() ([]int64, error) {
	var userIDs []int64
	err := r.db.Model(&models.Expense{}).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// SumByFilter returns the total spent amount for a user matching the given filter.
func (r *ExpenseRepository) SumByFilter(userID int64, filter models.ExpenseFilter) (float64, error) {
	query := r.db.Model(&models.Expense{}).Where("user_id = ?", userID)

	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("category = ?", *filter.Category)
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", *filter.DateFrom)
		if err == nil {
			query = query.Where("spent_at >= ?", dateFrom)
		}
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", *filter.DateTo)
		if err == nil {
			query = query.Where("spent_at <= ?", dateTo)
		}
	}

	var total float64
	err := query.Select("COALESCE(SUM(amount), 0)").Scan(&total).Error
	return total, err
}
