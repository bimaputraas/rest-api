package redis_repo

import (
	"context"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
	"github.com/bimaputraas/rest-api/internal/repository"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	"github.com/go-redis/redis"
)

func Cacher(client *redis.Client) repository.Cacher {
	return &cacher{client}
}

type (
	cacher struct {
		*redis.Client
	}
)

func (r *cacher) GetBalanceByUId(ctx context.Context, userId uint) (model.Balance, error) {
	return model.Balance{}, pkgerrors.Unimplemented()
}
func (r *cacher) GetTransactionsByUId(ctx context.Context, userId uint) ([]model.Transaction, error) {
	return nil, pkgerrors.Unimplemented()
}

// Set ttl to 0 for no expiration time.
func (r *cacher) InsertTransactions(ctx context.Context, transaction []model.Transaction, ttl time.Duration) error {
	return pkgerrors.Unimplemented()
}

// Set ttl to 0 for no expiration time.
func (r *cacher) InsertBalance(ctx context.Context, balance model.Balance, ttl time.Duration) error {
	return pkgerrors.Unimplemented()
}
