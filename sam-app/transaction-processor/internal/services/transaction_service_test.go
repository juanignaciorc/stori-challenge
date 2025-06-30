package services

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"transaction-processor/internal/domain/model"
	"transaction-processor/internal/mocks"
)

func TestTransactionService_ProcessTransactionsAndSendSummary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileReader := mocks.NewMockFileReader(ctrl)
	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	mockRepo := mocks.NewMockTransactionRepository(ctrl)

	service := NewTransactionService(mockFileReader, mockEmailSender, mockRepo)

	filePath := "transactions.csv"
	email := "user@example.com"
	accountID := "acc123"

	tx1 := &model.Transaction{ID: "tx1", Amount: 100, IsCredit: true, Date: time.Now()}
	tx2 := &model.Transaction{ID: "tx2", Amount: 50, IsCredit: false, Date: time.Now()}

	// Mock: leer archivo
	mockFileReader.
		EXPECT().
		ReadTransactions(filePath).
		Return([]*model.Transaction{tx1, tx2}, nil)

	// Mock: guardar transacciones
	mockRepo.EXPECT().SaveTransaction(tx1).Return(nil)
	mockRepo.EXPECT().SaveTransaction(tx2).Return(nil)

	// Mock: guardar resumen
	gomock.InOrder(
		mockRepo.EXPECT().SaveAccount(gomock.Any(), gomock.Any()).Return(nil),
		mockEmailSender.EXPECT().SendSummaryEmail(email, gomock.Any()).Return(nil),
	)

	err := service.ProcessTransactionsAndSendSummary(filePath, email, accountID)
	assert.NoError(t, err)
}
