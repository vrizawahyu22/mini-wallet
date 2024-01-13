package querier

import (
	"mini-wallet/utils"

	"github.com/jackc/pgx/v4"
)

type Repository interface {
	Querier

	WithTx(tx pgx.Tx) Querier
	GetDB() utils.PGXPool
}

type RepositoryImpl struct {
	db utils.PGXPool
	*Queries
}

func NewRepository(db utils.PGXPool) Repository {
	return &RepositoryImpl{db: db, Queries: New(db)}
}

func (r *RepositoryImpl) WithTx(tx pgx.Tx) Querier {
	return &Queries{
		db: tx,
	}
}

func (r *RepositoryImpl) GetDB() utils.PGXPool {
	return r.db
}
