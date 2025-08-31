package accounts

import (
	"github.com/ziflex/dbx"
)

type Repository interface {
	CreateAccount(ctx dbx.Context, acc AccountCreation) (Account, error)
	GetAccountByID(ctx dbx.Context, id int64) (Account, error)
}
