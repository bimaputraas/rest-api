package usecase

import (
	"context"
	"fmt"
	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkgvalidate "github.com/bimaputraas/rest-api/pkg/validate"
	"sync"
	"time"
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

	_, err := u.repo.GetUserById(ctx, targetUid)
	if pkgerrors.Code(err) == pkgerrors.ErrNotFound {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("target user is not exist"))
	}
	if err != nil {
		return model.Transfer{}, pkgerrors.Internal(err)
	}

	if amount < 1 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	uBalance, err := u.repo.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.Transfer{}, err
	}

	balanceBefore := uBalance.CurrentBalance
	balanceAfter := balanceBefore - amount
	if balanceAfter < 0 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("balance is not enough"))
	}

	uBalanceTarget, err := u.repo.GetBalanceByUId(ctx, targetUid)
	if err != nil {
		return model.Transfer{}, err
	}

	balanceBeforeTarget := uBalanceTarget.CurrentBalance
	balanceAfterTarget := balanceBeforeTarget + amount

	uBalance.CurrentBalance = balanceAfter
	uBalance.Updated = now

	txRepo, err := u.repo.BeginTx()
	if err != nil {
		return model.Transfer{}, err
	}

	wg.Add(3)
	go func() {
		defer wg.Done()
		err = txRepo.UpdateBalance(ctx, uBalanceTarget)
		if err != nil {
			mutex.Lock()
			errs = append(errs, err)
			mutex.Unlock()
		}

	}()

	uBalanceTarget.CurrentBalance = balanceAfterTarget
	uBalanceTarget.Updated = now

	go func() {
		defer wg.Done()
		err = txRepo.UpdateBalance(ctx, uBalanceTarget)
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
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
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
