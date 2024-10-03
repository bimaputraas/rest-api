package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkgvalidate "github.com/bimaputraas/rest-api/pkg/validate"
)

type (
	Transfer struct {
		TargetUser uint    `validate:"required" json:"target_user"`
		Amount     float64 `validate:"required" json:"amount"`
		Remarks    string  `validate:"required" json:"remarks"`
	}
)

func (u *Usecase) Transfer(ctx context.Context, userId uint, transfer Transfer) (model.Transfer, error) {

	var (
		data  = model.Transfer{}
		errs  = []error{}
		wg    = sync.WaitGroup{}
		mutex = sync.Mutex{}
	)
	if err := pkgvalidate.Struct(transfer); err != nil {
		return model.Transfer{}, pkgerrors.InvalidArgument(err)
	}

	remarks := transfer.Remarks
	amount := transfer.Amount
	targetUid := transfer.TargetUser
	now := time.Now().Format(time.DateTime)

	_, err := u.repo.Storage.GetUserById(ctx, targetUid)
	if pkgerrors.Code(err) == pkgerrors.ErrNotFound {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("target user is not exist"))
	}
	if err != nil {
		return model.Transfer{}, pkgerrors.Internal(err)
	}

	if amount < 1 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	userBalance, err := u.repo.Storage.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.Transfer{}, err
	}

	userBalanceBefore := userBalance.CurrentBalance
	userBalanceAfter := userBalanceBefore - amount
	if userBalanceAfter < 0 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("balance is not enough"))
	}

	targetBalance, err := u.repo.Storage.GetBalanceByUId(ctx, targetUid)
	if err != nil {
		return model.Transfer{}, err
	}

	targetBalanceBefore := targetBalance.CurrentBalance
	targetBalanceAfter := targetBalanceBefore + amount

	userBalance.CurrentBalance = userBalanceAfter
	userBalance.Updated = now

	txRepo, err := u.repo.Storage.BeginTx()
	if err != nil {
		return model.Transfer{}, err
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		err = txRepo.UpdateBalance(ctx, userBalance)
		if err != nil {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}

	}()

	targetBalance.CurrentBalance = targetBalanceAfter
	targetBalance.Updated = now

	go func() {
		defer wg.Done()
		err = txRepo.UpdateBalance(ctx, targetBalance)
		if err != nil {

			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}

	}()

	go func() {
		defer wg.Done()
		r, err := txRepo.InsertTransfer(ctx, model.Transfer{
			UserID:        userId,
			Amount:        amount,
			Remarks:       remarks,
			BalanceBefore: userBalanceBefore,
			BalanceAfter:  userBalanceAfter,
			Created:       now,
		})

		if err != nil {

			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}
		data = r

		transaction := model.Transaction{
			UserID:        userId,
			TransferId:    &data.ID,
			Amount:        data.Amount,
			Remarks:       data.Remarks,
			BalanceBefore: data.BalanceBefore,
			BalanceAfter:  data.BalanceAfter,
			Status:        "SUCCESS",
			Created:       now,
		}

		transaction.MockRandType()
		_, err = txRepo.InsertTransaction(ctx, transaction)
		if err != nil {

			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}

	}()

	wg.Wait()

	if len(errs) > 0 {
		errRollback := txRepo.Rollback()
		if err != nil {
			return model.Transfer{}, errRollback
		}
		return model.Transfer{}, errs[0]
	}

	data.UserID = 0

	return data, txRepo.Commit()
}
