package adapters

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"transaction-processor/internal/domain/model"
)

// CSVFileReader implements the FileReader port for CSV files
type CSVFileReader struct{}

// NewCSVFileReader creates a new CSVFileReader
func NewCSVFileReader() *CSVFileReader {
	return &CSVFileReader{}
}

// ReadTransactions reads transactions from a CSV file
func (r *CSVFileReader) ReadTransactions(filePath string) ([]*model.Transaction, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read the header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// Validate header
	if len(header) < 3 || header[0] != "Id" || header[1] != "Date" || header[2] != "Transaction" {
		return nil, fmt.Errorf("invalid CSV format, expected header: Id,Date,Transaction")
	}

	// Read transactions
	var transactions []*model.Transaction
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV record: %w", err)
		}

		// Validate record
		if len(record) < 3 {
			return nil, fmt.Errorf("invalid CSV record format, expected at least 3 fields")
		}

		// Create transaction
		tx, err := model.NewTransaction(record[0], record[1], record[2])
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %w", err)
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
