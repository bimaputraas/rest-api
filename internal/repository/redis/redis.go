package redis

import (
	"context"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
	"github.com/bimaputraas/rest-api/internal/repository"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	"github.com/go-redis/redis"
)

func NewCache(client *redis.Client) repository.Cache {
	return &cache{client}
}

type (
	cache struct {
		*redis.Client
	}
)

func (c *cache) GetBalanceByUId(ctx context.Context, userId uint) (model.Balance, error) {
	return model.Balance{}, pkgerrors.Unimplemented()
}
func (c *cache) GetTransactionsByUId(ctx context.Context, userId uint) ([]model.Transaction, error) {
	return nil, pkgerrors.Unimplemented()
}

// Set ttl to 0 for no expiration time.
func (c *cache) InsertTransactions(ctx context.Context, transaction []model.Transaction, ttl time.Duration) error {
	return pkgerrors.Unimplemented()
}

// Set ttl to 0 for no expiration time.
func (c *cache) InsertBalance(ctx context.Context, balance model.Balance, ttl time.Duration) error {
	return pkgerrors.Unimplemented()
}
