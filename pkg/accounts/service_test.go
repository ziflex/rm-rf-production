package accounts_test

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/ziflex/dbx"
	"github.com/ziflex/rm-rf-production/internal/database"
	"github.com/ziflex/rm-rf-production/pkg/accounts"
	"github.com/ziflex/rm-rf-production/pkg/common"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateAccount_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	expected := accounts.Account{
		ID:             1,
		DocumentNumber: "abc",
	}

	mock.ExpectBegin().WillReturnError(nil)
	mock.ExpectQuery(`INSERT INTO accounts \(document_number\) VALUES \(\$1\) RETURNING id`).
		WithArgs(expected.DocumentNumber).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(1),
		)
	mock.ExpectCommit()

	actual, err := svc.CreateAccount(context.Background(), accounts.AccountCreation{DocumentNumber: expected.DocumentNumber})

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_CreateAccount_Error_Duplicate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	mock.ExpectBegin().WillReturnError(nil)
	mock.ExpectQuery(`INSERT INTO accounts \(document_number\) VALUES \(\$1\) RETURNING id`).
		WithArgs("abc").
		WillReturnError(
			&pq.Error{
				Code: "23505",
			},
		)
	mock.ExpectRollback()

	_, err = svc.CreateAccount(context.Background(), accounts.AccountCreation{DocumentNumber: "abc"})

	assert.ErrorIs(t, err, common.ErrDuplicate)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_CreateAccount_Error_Propagated(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	mock.ExpectBegin().WillReturnError(nil)
	mock.ExpectQuery(`INSERT INTO accounts \(document_number\) VALUES \(\$1\) RETURNING id`).WillReturnError(
		&pq.Error{
			Code: "08006",
		},
	)
	mock.ExpectRollback()

	_, err = svc.CreateAccount(context.Background(), accounts.AccountCreation{DocumentNumber: "abc"})

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_GetAccountByID_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	mock.ExpectQuery(`SELECT \* FROM accounts WHERE id=\$1`).
		WithArgs(7).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "document_number"}).
				AddRow(7, "abc"),
		)

	expected := accounts.Account{
		ID:             7,
		DocumentNumber: "abc",
	}

	actual, err := svc.GetAccountByID(context.Background(), 7)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_GetAccountByID_Error_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	mock.ExpectQuery(`SELECT \* FROM accounts WHERE id=\$1`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "document_number"}),
	)

	_, err = svc.GetAccountByID(context.Background(), 7)

	assert.ErrorIs(t, err, common.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestService_GetAccountByID_Error_Propagated(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	db := dbx.New(mockDB)
	svc := accounts.NewService(db, database.NewAccountsRepository())

	mock.ExpectQuery(`SELECT \* FROM accounts WHERE id=\$1`).
		WithArgs(7).
		WillReturnError(
			&pq.Error{
				Code: "08006",
			},
		)

	_, err = svc.GetAccountByID(context.Background(), 7)

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
