package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/bimaputraas/rest-api/internal/model"
	"github.com/bimaputraas/rest-api/internal/repository"

	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	"gorm.io/gorm"
)

func NewStorage(gormDB *gorm.DB) repository.Storage {
	return &storage{
		DB: gormDB,
	}
}

type (
	storage struct {
		*gorm.DB
	}
)

func (s *storage) BeginTx() (repository.Tx, error) {
	return &storage{DB: s.DB.Begin()}, nil
}

func (s *storage) Commit() error {
	s.DB.Commit()
	return nil
}

func (s *storage) Rollback() error {
	s.DB.Rollback()
	return nil
}

func (s *storage) GetUserById(ctx context.Context, id uint) (model.User, error) {
	user := model.User{}
	result := s.DB.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, pkgerrors.NotFound(result.Error)
	}

	return user, result.Error
}

func (s *storage) GetUserByPhone(ctx context.Context, phone string) (model.User, error) {
	user := model.User{}
	result := s.DB.Where("phone_number = ?", phone).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model.User{}, pkgerrors.NotFound(result.Error)
	}

	return user, result.Error
}

func (s *storage) GetBalanceByUId(ctx context.Context, Uid uint) (model.Balance, error) {
	balance := model.Balance{}
	result := s.DB.Where("user_id = ?", Uid).First(&balance)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return balance, pkgerrors.NotFound(result.Error)
	}

	return balance, result.Error
}

func (s *storage) GetTransactionsByUId(ctx context.Context, id uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	result := s.DB.Where("user_id = ?", id).Find(&transactions)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return transactions, pkgerrors.NotFound(result.Error)
	}
	return transactions, result.Error
}

func (s *storage) InsertUser(ctx context.Context, user model.User) (model.User, error) {
	result := s.DB.Create(&user)
	return user, result.Error
}

func (s *storage) InsertBalance(ctx context.Context, balance model.Balance) error {
	fmt.Println("balance", balance)
	result := s.DB.Create(&balance)
	return result.Error
}

func (s *storage) UpdateBalance(ctx context.Context, update model.Balance) error {
	return s.DB.Save(&update).Error
}

func (s *storage) InsertTopUp(ctx context.Context, topUp model.TopUp) (model.TopUp, error) {
	result := s.DB.Create(&topUp)
	return topUp, result.Error
}

func (s *storage) InsertPayment(ctx context.Context, payment model.Payment) (model.Payment, error) {
	result := s.DB.Create(&payment)
	return payment, result.Error
}

func (s *storage) InsertTransfer(ctx context.Context, transfer model.Transfer) (model.Transfer, error) {
	result := s.DB.Create(&transfer)
	return transfer, result.Error
}

func (s *storage) InsertTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	result := s.DB.Create(&transaction)
	return transaction, result.Error
}

func (s *storage) InsertTransactions(ctx context.Context, transaction []model.Transaction) error {
	result := s.DB.Create(&transaction)
	return result.Error
}
