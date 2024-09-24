package repository

import (
	"context"
	"github.com/bimaputraas/rest-api/internal/model"
)

const (
	InvalidArgument = 1
	NotFound        = 2
	Internal        = 3
	Illegal         = 4
)

type (
	Repository interface {
		ReadWriter
		BeginTx() (Tx, error)
	}

	Tx interface {
		Writer
		Commit() error
		Rollback() error
	}

	Reader interface {
		GetUserById(ctx context.Context, id uint) (model.User, error)
		GetUserByPhone(ctx context.Context, phone string) (model.User, error)
		GetBalanceByUId(ctx context.Context, Uid uint) (model.Balance, error)
		GetTransactionsByUId(ctx context.Context, id uint) ([]model.Transaction, error)
	}

	Writer interface {
		InsertUser(ctx context.Context, user model.User) (model.User, error)
		InsertBalance(ctx context.Context, balance model.Balance) error
		InsertTopUp(ctx context.Context, topUp model.TopUp) (model.TopUp, error)
		InsertPayment(ctx context.Context, payment model.Payment) (model.Payment, error)
		InsertTransfer(ctx context.Context, transfer model.Transfer) (model.Transfer, error)
		InsertTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
		UpdateBalance(ctx context.Context, update model.Balance) error
	}

	ReadWriter interface {
		Reader
		Writer
	}
)
