package model

import (
	"time"
)

// MonthlyStats represents transaction statistics for a specific month
type MonthlyStats struct {
	Month            time.Month
	Year             int
	TransactionCount int
	TotalCredit      float64
	TotalDebit       float64
	CreditCount      int
	DebitCount       int
}

// Account represents a financial account with transactions
type Account struct {
	Transactions []*Transaction
	Balance      float64
	MonthlyStats map[string]*MonthlyStats // key is "YYYY-MM"
}

// NewAccount creates a new empty account
func NewAccount() *Account {
	return &Account{
		Transactions: []*Transaction{},
		Balance:      0,
		MonthlyStats: make(map[string]*MonthlyStats),
	}
}

// AddTransaction adds a transaction to the account and updates the balance and stats
func (a *Account) AddTransaction(tx *Transaction) {
	a.Transactions = append(a.Transactions, tx)

	// Update balance
	if tx.IsCredit {
		a.Balance += tx.Amount
	} else {
		a.Balance -= tx.Amount
	}

	// Update monthly stats
	monthKey := formatMonthKey(tx.Date)
	stats, exists := a.MonthlyStats[monthKey]
	if !exists {
		stats = &MonthlyStats{
			Month: tx.Date.Month(),
			Year:  tx.Date.Year(),
		}
		a.MonthlyStats[monthKey] = stats
	}

	stats.TransactionCount++

	if tx.IsCredit {
		stats.TotalCredit += tx.Amount
		stats.CreditCount++
	} else {
		stats.TotalDebit += tx.Amount
		stats.DebitCount++
	}
}

// GetTotalBalance returns the current account balance
func (a *Account) GetTotalBalance() float64 {
	return a.Balance
}

// GetMonthlyTransactionCounts returns a map of month names to transaction counts
func (a *Account) GetMonthlyTransactionCounts() map[string]int {
	result := make(map[string]int)

	for _, stats := range a.MonthlyStats {
		monthName := stats.Month.String()
		result[monthName] = stats.TransactionCount
	}

	return result
}

// GetAverageCreditAmount returns the average credit amount across all transactions
func (a *Account) GetAverageCreditAmount() float64 {
	var totalCredit float64
	var creditCount int

	for _, stats := range a.MonthlyStats {
		totalCredit += stats.TotalCredit
		creditCount += stats.CreditCount
	}

	if creditCount == 0 {
		return 0
	}

	return totalCredit / float64(creditCount)
}

// GetAverageDebitAmount returns the average debit amount across all transactions
func (a *Account) GetAverageDebitAmount() float64 {
	var totalDebit float64
	var debitCount int

	for _, stats := range a.MonthlyStats {
		totalDebit += stats.TotalDebit
		debitCount += stats.DebitCount
	}

	if debitCount == 0 {
		return 0
	}

	return totalDebit / float64(debitCount)
}

// formatMonthKey formats a date as "YYYY-MM" for use as a map key
func formatMonthKey(date time.Time) string {
	return date.Format("2006-01")
}
