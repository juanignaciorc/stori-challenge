package adapters_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"transaction-processor/internal/adapters"
)

func TestCSVFileReader_ReadTransactions(t *testing.T) {
	// Create a temporary directory for the test file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test_transactions.csv")

	// Define the CSV content with 4 data rows
	csvContent := `Id,Date,Transaction
0,7/15,+60.5
1,7/28,-10.3
2,8/2,-20.46
3,8/13,+10`

	// Write the content to the temporary file
	err := os.WriteFile(testFilePath, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test CSV file: %v", err)
	}

	// Create an instance of CSVFileReader
	reader := adapters.NewCSVFileReader()

	// Call the method under test
	transactions, err := reader.ReadTransactions(testFilePath)
	if err != nil {
		t.Fatalf("ReadTransactions failed: %v", err)
	}

	// Get the current year to correctly format expected dates
	currentYear := time.Now().Year()

	// Assertions for the number of transactions
	if len(transactions) != 4 {
		t.Errorf("Expected 4 transactions, got %d", len(transactions))
	}

	// Check first transaction: +60.5 (Credit)
	expectedDate1 := time.Date(currentYear, time.July, 15, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	if transactions[0].ID != "0" || transactions[0].Date.Format("2006-01-02") != expectedDate1 || transactions[0].Amount != 60.5 || !transactions[0].IsCredit {
		t.Errorf("First transaction mismatch: Expected ID '0', Date '%s', Amount '60.5', IsCredit 'true', got %+v", expectedDate1, transactions[0])
	}

	// Check second transaction: -10.3 (Debit)
	expectedDate2 := time.Date(currentYear, time.July, 28, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	if transactions[1].ID != "1" || transactions[1].Date.Format("2006-01-02") != expectedDate2 || transactions[1].Amount != 10.3 || transactions[1].IsCredit {
		t.Errorf("Second transaction mismatch: Expected ID '1', Date '%s', Amount '10.3', IsCredit 'false', got %+v", expectedDate2, transactions[1])
	}

	// Check third transaction: -20.46 (Debit)
	expectedDate3 := time.Date(currentYear, time.August, 2, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	if transactions[2].ID != "2" || transactions[2].Date.Format("2006-01-02") != expectedDate3 || transactions[2].Amount != 20.46 || transactions[2].IsCredit {
		t.Errorf("Third transaction mismatch: Expected ID '2', Date '%s', Amount '20.46', IsCredit 'false', got %+v", expectedDate3, transactions[2])
	}

	// Check fourth transaction: +10 (Credit)
	expectedDate4 := time.Date(currentYear, time.August, 13, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	if transactions[3].ID != "3" || transactions[3].Date.Format("2006-01-02") != expectedDate4 || transactions[3].Amount != 10.0 || !transactions[3].IsCredit {
		t.Errorf("Fourth transaction mismatch: Expected ID '3', Date '%s', Amount '10.0', IsCredit 'true', got %+v", expectedDate4, transactions[3])
	}

	// --- Test Case: Invalid header ---
	t.Run("invalid header", func(t *testing.T) {
		invalidCsvContent := `WrongId,Date,Transaction`
		invalidFilePath := filepath.Join(tempDir, "invalid_header.csv")
		os.WriteFile(invalidFilePath, []byte(invalidCsvContent), 0644)

		_, err := reader.ReadTransactions(invalidFilePath)
		if err == nil {
			t.Error("Expected an error for invalid header, got nil")
		}
	})

	// --- Test Case: Non-existent file ---
	t.Run("non-existent file", func(t *testing.T) {
		_, err := reader.ReadTransactions(filepath.Join(tempDir, "non_existent.csv"))
		if err == nil {
			t.Error("Expected an error for non-existent file, got nil")
		}
	})
}
