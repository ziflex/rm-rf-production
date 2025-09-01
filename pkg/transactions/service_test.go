package transactions_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/ziflex/dbx"
	"github.com/ziflex/rm-rf-production/internal/database"
	"github.com/ziflex/rm-rf-production/pkg/common"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
)

func TestService_CreateTransaction_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := transactions.NewService(db, database.NewTransactions())

	type testCase struct {
		Name          string
		OperationType transactions.OperationType
		AmountIn      float64
		AmountOut     float64
	}

	tsdata := []testCase{
		{"Purchase", transactions.OperationTypePurchase, 123.45, -123.45},
		{"InstallmentPurchase", transactions.OperationTypeInstallmentPurchase, 67.89, -67.89},
		{"Withdrawal", transactions.OperationTypeWithdrawal, 10.00, -10.00},
		{"Payment", transactions.OperationTypePayment, 200.00, 200.00},
	}

	var txId int64 = 1
	var txAccountId int64 = 5

	for _, tc := range tsdata {
		t.Run(tc.Name, func(t *testing.T) {
			ts := time.Now()

			mock.ExpectBegin().WillReturnError(nil)
			mock.ExpectQuery(
				`INSERT INTO transactions \(account_id, operation_type, amount\) VALUES \(\$1, \$2, \$3\) RETURNING id, account_id, operation_type, amount, event_date`,
			).
				WithArgs(txAccountId, tc.OperationType.String(), tc.AmountOut).
				WillReturnRows(sqlmock.
					NewRows([]string{"id", "account_id", "operation_type", "amount", "event_date"}).
					AddRow(txId, txAccountId, tc.OperationType.String(), tc.AmountOut, ts),
				)
			mock.ExpectCommit()

			expected := transactions.Transaction{
				ID:            txId,
				AccountID:     txAccountId,
				OperationType: tc.OperationType,
				Amount:        tc.AmountOut,
				EventDate:     ts,
			}

			actual, err := svc.CreateTransaction(context.Background(), transactions.TransactionCreation{
				AccountID:     txAccountId,
				OperationType: tc.OperationType,
				Amount:        tc.AmountIn,
			})

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestService_CreateTransaction_Error_InvalidAmount(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := transactions.NewService(db, database.NewTransactions())

	_, err = svc.CreateTransaction(context.Background(), transactions.TransactionCreation{
		AccountID:     100,
		OperationType: transactions.OperationTypePurchase,
		Amount:        0,
	})

	assert.ErrorIs(t, err, transactions.ErrInvalidAmount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_CreateTransaction_Error_InvalidOperation(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := transactions.NewService(db, database.NewTransactions())

	_, err = svc.CreateTransaction(context.Background(), transactions.TransactionCreation{
		AccountID:     100,
		OperationType: transactions.OperationType(999), // Invalid operation type
		Amount:        10000,
	})

	assert.ErrorIs(t, err, transactions.ErrInvalidOperationType)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_CreateTransaction_Error_InvalidAccountID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := transactions.NewService(db, database.NewTransactions())

	var accId int64 = 0 // Invalid account ID
	var opType transactions.OperationType = transactions.OperationTypePurchase
	var amt float64 = 123.45

	mock.ExpectBegin().WillReturnError(nil)
	mock.ExpectQuery(`.*`).
		WithArgs(accId, opType.String(), -amt).
		WillReturnError(
			&pq.Error{
				Code: "23503",
			},
		)
	mock.ExpectRollback()

	_, err = svc.CreateTransaction(context.Background(), transactions.TransactionCreation{
		AccountID:     accId,
		OperationType: opType,
		Amount:        amt,
	})

	assert.ErrorIs(t, err, common.ErrNotFound)
}

func TestService_CreateTransaction_Error_Propagate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := transactions.NewService(db, database.NewTransactions())

	var accId int64 = 0 // Invalid account ID
	var opType transactions.OperationType = transactions.OperationTypePurchase
	var amt float64 = 123.45

	mock.ExpectBegin().WillReturnError(nil)
	mock.ExpectQuery(`.*`).
		WithArgs(accId, opType.String(), -amt).
		WillReturnError(
			&pq.Error{
				Code: "22004",
			},
		)
	mock.ExpectRollback()

	_, err = svc.CreateTransaction(context.Background(), transactions.TransactionCreation{
		AccountID:     accId,
		OperationType: opType,
		Amount:        amt,
	})

	assert.Error(t, err)
}
