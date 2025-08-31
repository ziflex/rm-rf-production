package api

import (
	"context"

	"github.com/ziflex/rm-rf-production/pkg/accounts"
	"github.com/ziflex/rm-rf-production/pkg/transactions"
)

type Handler struct {
	accounts     accounts.Service
	transactions transactions.Service
}

func NewHandler(
	accounts accounts.Service,
	transactions transactions.Service,
) StrictServerInterface {
	return &Handler{
		accounts,
		transactions,
	}
}

func (r *Handler) CreateAccount(ctx context.Context, request CreateAccountRequestObject) (CreateAccountResponseObject, error) {
	acc, err := r.accounts.CreateAccount(ctx, accounts.AccountCreation{
		DocumentNumber: request.Body.DocumentNumber,
	})

	if err != nil {
		return nil, err
	}

	return CreateAccount201JSONResponse{
		AccountId:      acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}, nil
}

func (r *Handler) GetAccount(ctx context.Context, request GetAccountRequestObject) (GetAccountResponseObject, error) {
	acc, err := r.accounts.GetAccountByID(ctx, request.AccountId)

	if err != nil {
		return nil, err
	}

	return GetAccount200JSONResponse{
		AccountId:      acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}, nil
}

func (r *Handler) CreateTransaction(ctx context.Context, request CreateTransactionRequestObject) (CreateTransactionResponseObject, error) {
	tx, err := r.transactions.CreateTransaction(ctx, transactions.TransactionCreation{
		AccountID:     request.Body.AccountId,
		OperationType: transactions.NewOperationType(int(request.Body.OperationTypeId)),
		Amount:        request.Body.Amount,
	})

	if err != nil {
		return nil, err
	}

	return CreateTransaction201JSONResponse{
		TransactionId:   tx.ID,
		AccountId:       tx.AccountID,
		OperationTypeId: OperationType(tx.OperationType),
		Amount:          tx.Amount,
		EventDate:       tx.EventDate,
	}, nil
}
