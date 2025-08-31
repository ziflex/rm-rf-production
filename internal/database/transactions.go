package database

import (
	"database/sql"
	"fmt"

	"github.com/ziflex/dbx"
	"github.com/ziflex/rm-rf-production/pkg/common"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
)

type TransactionsRepository struct {
}

func NewTransactions() transactions.Repository {
	return &TransactionsRepository{}
}

func (t *TransactionsRepository) CreateTransaction(ctx dbx.Context, tr transactions.TransactionCreation) (transactions.Transaction, error) {
	row := ctx.Executor().QueryRow(`
		INSERT INTO transactions (account_id, operation_type, amount) VALUES ($1, $2, $3)
		RETURNING id, account_id, operation_type, amount, event_date
	`, tr.AccountID, tr.OperationType.String(), tr.Amount)

	if err := row.Err(); err != nil {
		if pgErr, ok := IsPgErr(err); ok {
			if IsDbForeignKeyViolation(pgErr) {
				return transactions.Transaction{}, fmt.Errorf("account %w: %d", common.ErrNotFound, tr.AccountID)
			}
		}

		return transactions.Transaction{}, err
	}

	res, err := t.scanTransaction(row)

	if err != nil {
		return transactions.Transaction{}, err
	}

	return res, nil
}

func (t *TransactionsRepository) scanTransaction(row *sql.Row) (transactions.Transaction, error) {
	var tr transactions.Transaction
	var optype string

	err := row.Scan(&tr.ID, &tr.AccountID, &optype, &tr.Amount, &tr.EventDate)

	if err != nil {
		return transactions.Transaction{}, err
	}

	tr.OperationType = transactions.NewOperationTypeFromString(optype)

	return tr, nil
}
