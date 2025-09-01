package database

import (
	"database/sql"
	"fmt"

	"github.com/ziflex/dbx"
	"github.com/ziflex/rm-rf-production/pkg/accounts"
	"github.com/ziflex/rm-rf-production/pkg/common"
)

type Accounts struct {
}

func NewAccountsRepository() accounts.Repository {
	return &Accounts{}
}

func (a *Accounts) CreateAccount(ctx dbx.Context, acc accounts.AccountCreation) (accounts.Account, error) {
	row := ctx.Executor().QueryRow(`
		INSERT INTO accounts (document_number) VALUES ($1)
		RETURNING id
	`, acc.DocumentNumber)

	if err := row.Err(); err != nil {
		if pgErr, ok := IsPgErr(err); ok {
			if IsDbUniqueViolation(pgErr) {
				return accounts.Account{}, fmt.Errorf("document number %w: %s", common.ErrDuplicate, acc.DocumentNumber)
			}
		}

		return accounts.Account{}, err
	}

	var id int64
	err := row.Scan(&id)

	if err != nil {
		return accounts.Account{}, err
	}

	return accounts.Account{
		ID:             id,
		DocumentNumber: acc.DocumentNumber,
	}, nil
}

func (a *Accounts) GetAccountByID(ctx dbx.Context, id int64) (accounts.Account, error) {
	rows, err := ctx.Executor().Query("SELECT * FROM accounts WHERE id=$1", id)

	if err != nil {
		return accounts.Account{}, err
	}

	defer rows.Close()

	if !rows.Next() {
		return accounts.Account{}, fmt.Errorf("account %w: %d", common.ErrNotFound, id)
	}

	return a.scanAccount(rows)
}

func (a *Accounts) scanAccount(rows *sql.Rows) (accounts.Account, error) {
	var acc accounts.Account
	err := rows.Scan(&acc.ID, &acc.DocumentNumber)

	if err != nil {
		return accounts.Account{}, err
	}

	return acc, nil
}
