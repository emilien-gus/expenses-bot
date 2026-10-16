package models

import (
	"time"
)

// Category represents a valid expense category.
type Category string

const (
	CategoryFood          Category = "еда"
	CategoryTransport     Category = "транспорт"
	CategoryHousing       Category = "жильё"
	CategoryEntertainment Category = "развлечения"
	CategoryHealth        Category = "здоровье"
	CategoryShopping      Category = "покупки"
	CategoryOther         Category = "другое"
)

// ValidCategories contains all allowed expense categories.
var ValidCategories = []Category{
	CategoryFood,
	CategoryTransport,
	CategoryHousing,
	CategoryEntertainment,
	CategoryHealth,
	CategoryShopping,
	CategoryOther,
}

// IsValidCategory checks if a given string is a valid category.
func IsValidCategory(cat string) bool {
	for _, c := range ValidCategories {
		if string(c) == cat {
			return true
		}
	}
	return false
}

// Expense represents a single expense record in the database.
type Expense struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int64     `json:"user_id" gorm:"index;not null"`
	Amount    float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	Category  string    `json:"category" gorm:"type:varchar(50);index;not null"`
	SpentAt   time.Time `json:"spent_at" gorm:"type:date;index;not null"`
	Comment   string    `json:"comment" gorm:"type:text;default:''"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExpenseFilter holds query parameters for filtering expenses.
type ExpenseFilter struct {
	UserID   *int64  `form:"user_id"`
	Category *string `form:"category"`
	DateFrom *string `form:"date_from"`
	DateTo   *string `form:"date_to"`
}

// CreateExpenseRequest represents the payload for creating an expense.
type CreateExpenseRequest struct {
	UserID   int64   `json:"user_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Category string  `json:"category" binding:"required"`
	SpentAt  string  `json:"spent_at" binding:"required"`
	Comment  string  `json:"comment"`
}

// UpdateExpenseRequest represents the payload for updating an expense.
type UpdateExpenseRequest struct {
	UserID   *int64   `json:"user_id"`
	Amount   *float64 `json:"amount" binding:"omitempty,gt=0"`
	Category *string  `json:"category"`
	SpentAt  *string  `json:"spent_at"`
	Comment  *string  `json:"comment"`
}
