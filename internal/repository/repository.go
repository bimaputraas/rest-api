package repository

import (
	"context"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
)

const (
	InvalidArgument = 1
	NotFound        = 2
	Internal        = 3
	Illegal         = 4
)

func New(db Storage, cache Cache) *Repository {
	return &Repository{
		Storage: db,
		Cache:   cache,
	}
}

type (
	Repository struct {
		Storage Storage
		Cache   Cache
	}

	Storage interface {
		StorageWriter
		GetUserById(ctx context.Context, userId uint) (model.User, error)
		GetUserByPhone(ctx context.Context, phone string) (model.User, error)
		GetBalanceByUId(ctx context.Context, userId uint) (model.Balance, error)
		GetTransactionsByUId(ctx context.Context, userId uint) ([]model.Transaction, error)
		BeginTx() (Tx, error)
		InsertTransactions(ctx context.Context, transaction []model.Transaction) error
	}
	Cache interface {
		GetBalanceByUId(ctx context.Context, userId uint) (model.Balance, error)
		GetTransactionsByUId(ctx context.Context, userId uint) ([]model.Transaction, error)

		//Set ttl to 0 for no expiration time.
		InsertTransactions(ctx context.Context, transaction []model.Transaction, ttl time.Duration) error
		//Set ttl to 0 for no expiration time.
		InsertBalance(ctx context.Context, balance model.Balance, ttl time.Duration) error
	}
	Tx interface {
		StorageWriter
		Commit() error
		Rollback() error
	}

	StorageWriter interface {
		InsertUser(ctx context.Context, user model.User) (model.User, error)
		InsertBalance(ctx context.Context, balance model.Balance) error
		InsertTopUp(ctx context.Context, topUp model.TopUp) (model.TopUp, error)
		InsertPayment(ctx context.Context, payment model.Payment) (model.Payment, error)
		InsertTransfer(ctx context.Context, transfer model.Transfer) (model.Transfer, error)
		InsertTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
		UpdateBalance(ctx context.Context, update model.Balance) error
	}
)
