package transactions

import (
	"github.com/ziflex/dbx"
)

type Repository interface {
	CreateTransaction(ctx dbx.Context, tr TransactionCreation) (Transaction, error)
}
