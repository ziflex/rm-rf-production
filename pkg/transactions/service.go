package transactions

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ziflex/dbx"
)

type (
	Service interface {
		CreateTransaction(ctx context.Context, creation TransactionCreation) (Transaction, error)
	}

	serviceImpl struct {
		db         dbx.Database
		repository Repository
	}
)

func NewService(
	db dbx.Database,
	repository Repository,
) Service {
	return &serviceImpl{
		db:         db,
		repository: repository,
	}
}

func (s *serviceImpl) CreateTransaction(ctx context.Context, creation TransactionCreation) (Transaction, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("creating transaction")

	amt, err := s.handleOperation(creation.OperationType, creation.Amount)

	if err != nil {
		log.Error().Err(err).Msg("failed to handle operation")
		return Transaction{}, err
	}

	return dbx.TransactionWithResult[Transaction](ctx, s.db, func(tx dbx.Context) (Transaction, error) {
		t, err := s.repository.CreateTransaction(tx, TransactionCreation{
			AccountID:     creation.AccountID,
			OperationType: creation.OperationType,
			Amount:        amt,
		})

		if err != nil {
			log.Error().Err(err).Msg("failed to create transaction")

			return Transaction{}, err
		}

		log.Info().Int64("transaction_id", t.ID).Msg("transaction created")

		return t, nil
	})
}

func (s *serviceImpl) handleOperation(op OperationType, amount float64) (float64, error) {
	switch op {
	case OperationTypePurchase, OperationTypeInstallmentPurchase, OperationTypeWithdrawal:
		return -amount, nil
	case OperationTypePayment:
		return +amount, nil
	default:
		return 0, ErrInvalidOperationType
	}
}
