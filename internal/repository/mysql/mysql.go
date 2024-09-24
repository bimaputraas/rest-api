package mysql

import (
	"context"
	"errors"
	"github.com/bimaputraas/rest-api/internal/model"
	"github.com/bimaputraas/rest-api/internal/repository"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	"gorm.io/gorm"
)

func New(db *gorm.DB) repository.Repository {
	return &repo{db}
}

type (
	repo struct {
		*gorm.DB
	}
)

func (r *repo) BeginTx() (repository.Tx, error) {
	return &repo{r.DB.Begin()}, nil
}

func (r *repo) Commit() error {
	r.DB.Commit()
	return nil
}

func (r *repo) Rollback() error {
	r.DB.Rollback()
	return nil
}

func (r *repo) GetUserById(ctx context.Context, id uint) (model.User, error) {
	user := model.User{}
	result := r.DB.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, pkgerrors.NotFound(result.Error)
	}

	return user, result.Error
}

func (r *repo) GetUserByPhone(ctx context.Context, phone string) (model.User, error) {
	user := model.User{}
	result := r.DB.Where("phone_number = ?", phone).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, pkgerrors.NotFound(result.Error)
	}

	return user, result.Error
}

func (r *repo) GetBalanceByUId(ctx context.Context, Uid uint) (model.Balance, error) {
	balance := model.Balance{}
	result := r.DB.Where("user_id = ?", Uid).First(&balance)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return balance, pkgerrors.NotFound(result.Error)
	}

	return balance, result.Error
}

func (r *repo) GetTransactionsByUId(ctx context.Context, id uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	result := r.DB.Where("user_id = ?", id).Find(&transactions)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return transactions, pkgerrors.NotFound(result.Error)
	}
	return transactions, result.Error
}

func (r *repo) InsertUser(ctx context.Context, user model.User) (model.User, error) {
	result := r.DB.Create(&user)
	return user, result.Error
}

func (r *repo) InsertBalance(ctx context.Context, balance model.Balance) error {
	result := r.DB.Create(&balance)
	return result.Error
}

func (r *repo) UpdateBalance(ctx context.Context, update model.Balance) error {
	return r.DB.Save(&update).Error
}

func (r *repo) InsertTopUp(ctx context.Context, topUp model.TopUp) (model.TopUp, error) {
	result := r.DB.Create(&topUp)
	return topUp, result.Error
}

func (r *repo) InsertPayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	result := r.DB.Create(&payment)
	return payment, result.Error
}

func (r *repo) InsertTransfer(ctx context.Context, transfer model.Transfer) (model.Transfer, error) {
	result := r.DB.Create(&transfer)
	return transfer, result.Error
}

func (r *repo) InsertTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	result := r.DB.Create(&transaction)
	return transaction, result.Error
}
