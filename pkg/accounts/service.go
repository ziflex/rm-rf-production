package accounts

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ziflex/dbx"
)

type Service struct {
	db         dbx.Database
	repository Repository
}

func NewService(db dbx.Database, repository Repository) *Service {
	return &Service{db, repository}
}

func (s *Service) CreateAccount(ctx context.Context, creation AccountCreation) (Account, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("creating account")

	return dbx.TransactionWithResult[Account](ctx, s.db, func(tx dbx.Context) (Account, error) {
		acc, err := s.repository.CreateAccount(tx, creation)

		if err != nil {
			log.Error().Err(err).Msg("failed to create account")

			return Account{}, err
		}

		log.Info().Int64("id", acc.ID).Msg("account created")

		return acc, nil
	})
}

func (s *Service) GetAccountByID(ctx context.Context, id int64) (Account, error) {
	log := zerolog.Ctx(ctx)
	log.Info().Int64("id", id).Msg("getting account")

	acc, err := s.repository.GetAccountByID(dbx.NewContextFrom(ctx, s.db), id)

	if err != nil {
		log.Error().Err(err).Int64("id", id).Msg("failed to get account")

		return Account{}, err
	}

	log.Info().Int64("id", acc.ID).Msg("account retrieved")

	return acc, nil
}
