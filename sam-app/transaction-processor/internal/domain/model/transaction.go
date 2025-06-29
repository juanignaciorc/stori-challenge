package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID       string
	Date     time.Time
	Amount   float64
	IsCredit bool
}

// NewTransaction creates a new Transaction from raw data
func NewTransaction(id, dateStr, amountStr string) (*Transaction, error) {
	// Parse amount
	amount, err := parseAmount(amountStr)
	if err != nil {
		return nil, err
	}

	// Determine if credit or debit
	isCredit := !strings.HasPrefix(amountStr, "-")

	// Parse date (assuming MM/DD format)
	date, err := parseDate(dateStr)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		ID:       id,
		Date:     date,
		Amount:   amount,
		IsCredit: isCredit,
	}, nil
}

// parseAmount parses the transaction amount string
func parseAmount(amountStr string) (float64, error) {
	// Remove + or - prefix for parsing
	cleanAmount := strings.TrimPrefix(strings.TrimPrefix(amountStr, "+"), "-")
	return strconv.ParseFloat(cleanAmount, 64)
}

// parseDate parses the date string in MM/DD format
// Assumes current year for simplicity
func parseDate(dateStr string) (time.Time, error) {
	parts := strings.Split(dateStr, "/")
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}

	// Use current year
	year := time.Now().Year()

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}
