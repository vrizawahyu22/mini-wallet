package service

import (
	"context"
	"fmt"
	"mini-wallet/constant"
	"mini-wallet/dto"
	"mini-wallet/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	querier "mini-wallet/db/repository"
)

type WalletSvc interface {
	InitializeAccount(ctx context.Context, customerXID string) dto.InitializeAccountResp
	EnableWallet(ctx context.Context) dto.WalletResp
	DisableWallet(ctx context.Context, isDisabled bool) dto.WalletDisabledResp
	GetWallet(ctx context.Context) dto.WalletResp
	GetTransactions(ctx context.Context) dto.TransactionResp
	DepositWallet(ctx context.Context, req dto.DepositReq) dto.DepositResp
	WithdrawalWallet(ctx context.Context, req dto.WithdrawalReq) dto.WithdrawalResp
}

type WalletSvcImpl struct {
	repo querier.Repository
}

func NewWalletSvc(
	repo querier.Repository,
) WalletSvc {
	return &WalletSvcImpl{
		repo: repo,
	}
}

func (s *WalletSvcImpl) InitializeAccount(ctx context.Context, customerXID string) dto.InitializeAccountResp {
	resp, err := s.repo.FindUser(ctx, customerXID)
	if err != nil {
		err = utils.CustomErrorWithTrace(err, "Failed to find customer", 400)
		utils.PanicIfError(err)
	}

	return dto.InitializeAccountResp{
		Token: resp.Token,
	}
}

func getUser(ctx context.Context, repo querier.Repository) string {
	authPayload := utils.GetRequestCtx(ctx, utils.UserSession)

	user, err := repo.FindTokenByUserId(ctx, authPayload.Token)
	if err != nil || user.ID == "" {
		err = utils.CustomError("unauthorized", 400)
		utils.PanicIfError(err)
	}

	return user.ID
}

func (s *WalletSvcImpl) EnableWallet(ctx context.Context) dto.WalletResp {
	userID := getUser(ctx, s.repo)

	isExists, err := s.repo.CheckWalletExists(ctx, userID)
	if err != nil {
		err = utils.CustomError("Failed to check wallet", 400)
		utils.PanicIfError(err)
	}

	var wallet querier.Wallet
	if !isExists {
		wallet, err = s.repo.CreateWallet(ctx, querier.CreateWalletParams{
			ID:      uuid.NewString(),
			UserID:  userID,
			Status:  constant.WALLET_STATUS_ENABLED,
			Balance: 0,
		})
		if err != nil {
			err = utils.CustomError("Failed to enable wallet", 400)
			utils.PanicIfError(err)
		}
	} else {
		wallet, err = s.repo.FindWalletByUserId(ctx, userID)
		if err != nil {
			err = utils.CustomError("Failed to find wallet", 400)
			utils.PanicIfError(err)
		}

		if wallet.Status == constant.WALLET_STATUS_ENABLED {
			err = utils.CustomError("Already Enabled", 400)
			utils.PanicIfError(err)
		}

		wallet, err = s.repo.UpdateWalletByUserId(ctx, querier.UpdateWalletByUserIdParams{
			Status: constant.WALLET_STATUS_ENABLED,
			UserID: userID,
		})
		if err != nil {
			err = utils.CustomError("Failed to enable wallet", 400)
			utils.PanicIfError(err)
		}
	}

	return dto.WalletResp{
		Wallet: dto.Wallet{
			ID:        wallet.ID,
			OwnedBy:   wallet.UserID,
			Status:    wallet.Status,
			Balance:   int64(wallet.Balance),
			EnabledAt: wallet.EnabledAt,
		},
	}
}

func (s *WalletSvcImpl) DisableWallet(ctx context.Context, isDisable bool) dto.WalletDisabledResp {
	userID := getUser(ctx, s.repo)

	isExists, err := s.repo.CheckWalletExists(ctx, userID)
	if err != nil || !isExists {
		err = utils.CustomError("Failed to check wallet", 400)
		utils.PanicIfError(err)
	}

	wallet, err := s.repo.FindWalletByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomError("Failed to find wallet", 400)
		utils.PanicIfError(err)
	}

	if wallet.Status == constant.WALLET_STATUS_DISABLED {
		err = utils.CustomError("Already Disabled", 400)
		utils.PanicIfError(err)
	}

	wallet, err = s.repo.UpdateWalletByUserId(ctx, querier.UpdateWalletByUserIdParams{
		Status: constant.WALLET_STATUS_DISABLED,
		UserID: userID,
	})
	if err != nil {
		err = utils.CustomError("Failed to enable wallet", 400)
		utils.PanicIfError(err)
	}

	return dto.WalletDisabledResp{
		Wallet: dto.WalletDisabled{
			ID:         wallet.ID,
			OwnedBy:    wallet.UserID,
			Status:     wallet.Status,
			Balance:    int64(wallet.Balance),
			DisabledAt: wallet.EnabledAt,
		},
	}
}

func (s *WalletSvcImpl) GetWallet(ctx context.Context) dto.WalletResp {
	userID := getUser(ctx, s.repo)

	wallet, err := s.repo.FindWalletByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomError("Failed to find wallet", 400)
		utils.PanicIfError(err)
	}

	if wallet.Status == constant.WALLET_STATUS_DISABLED {
		err = utils.CustomError("Wallet disabled", 400)
		utils.PanicIfError(err)
	}

	return dto.WalletResp{
		Wallet: dto.Wallet{
			ID:        wallet.ID,
			OwnedBy:   wallet.UserID,
			Status:    wallet.Status,
			Balance:   int64(wallet.Balance),
			EnabledAt: wallet.EnabledAt,
		},
	}
}

func (s *WalletSvcImpl) GetTransactions(ctx context.Context) dto.TransactionResp {
	userID := getUser(ctx, s.repo)

	wallet, err := s.repo.FindWalletByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomErrorWithTrace(err, "Failed to find wallet", 400)
		utils.PanicIfError(err)
	}

	if wallet.Status == constant.WALLET_STATUS_DISABLED {
		err = utils.CustomError("Wallet disabled", 400)
		utils.PanicIfError(err)
	}

	transaction, err := s.repo.FindTransactionByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomErrorWithTrace(err, "Failed to find transaction", 400)
		utils.PanicIfError(err)
	}

	resp := utils.TransformDataOrPanic(transaction, []dto.Transaction{})

	return dto.TransactionResp{
		Transaction: resp,
	}
}

func updateBalance(repo querier.Repository, balance float64, userID string) {
	time.Sleep(3 * time.Second)
	ctx := context.Background()
	fmt.Println("MULAI UPDATE BALANCE")

	_, err := repo.UpdatBalanceWalletByUserId(ctx, querier.UpdatBalanceWalletByUserIdParams{
		Balance: balance,
		UserID:  userID,
	})

	if err != nil {
		err = utils.CustomError("Failed to update balance", 400)
		utils.PanicIfError(err)
	}
}

func (s *WalletSvcImpl) DepositWallet(ctx context.Context, req dto.DepositReq) dto.DepositResp {
	userID := getUser(ctx, s.repo)

	wallet, err := s.repo.FindWalletByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomError("Failed to find wallet", 400)
		utils.PanicIfError(err)
	}

	if wallet.Status == constant.WALLET_STATUS_DISABLED {
		err = utils.CustomError("Wallet disabled", 400)
		utils.PanicIfError(err)
	}

	var transaction querier.Transaction
	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		transaction, err = repoTx.CreateTransaction(ctx, querier.CreateTransactionParams{
			ID:          uuid.NewString(),
			UserID:      userID,
			Status:      constant.TRANSACTION_STATUS_SUCCESS,
			Type:        "deposit",
			Balance:     req.Amount,
			ReferenceID: req.ReferenceID,
		})

		return err
	})

	if err != nil {
		err = utils.CustomErrorWithTrace(err, "reference_id must be unique", 400)
		utils.PanicIfError(err)
	}

	balanceUpdated := req.Amount + wallet.Balance
	go updateBalance(s.repo, balanceUpdated, userID)

	return dto.DepositResp{
		Deposit: dto.Deposit{
			ID:          transaction.ID,
			DepositBy:   userID,
			Status:      transaction.Status,
			DepositedAt: transaction.TransactionAt,
			Amount:      int64(transaction.Balance),
			ReferenceID: transaction.ReferenceID,
		},
	}
}

func (s *WalletSvcImpl) WithdrawalWallet(ctx context.Context, req dto.WithdrawalReq) dto.WithdrawalResp {
	userID := getUser(ctx, s.repo)

	wallet, err := s.repo.FindWalletByUserId(ctx, userID)
	if err != nil {
		err = utils.CustomError("Failed to find wallet", 400)
		utils.PanicIfError(err)
	}

	if wallet.Status == constant.WALLET_STATUS_DISABLED {
		err = utils.CustomError("Wallet disabled", 400)
		utils.PanicIfError(err)
	}

	var transaction querier.Transaction
	balanceUpdated := wallet.Balance - req.Amount
	err = utils.ExecTxPool(ctx, s.repo.GetDB(), func(tx pgx.Tx) error {
		repoTx := s.repo.WithTx(tx)

		status := constant.TRANSACTION_STATUS_SUCCESS
		if balanceUpdated < 0 {
			status = constant.TRANSACTION_STATUS_FAILED
		}

		transaction, err = repoTx.CreateTransaction(ctx, querier.CreateTransactionParams{
			ID:          uuid.NewString(),
			UserID:      userID,
			Status:      status,
			Type:        "withdrawal",
			Balance:     req.Amount,
			ReferenceID: req.ReferenceID,
		})

		return err
	})

	if err != nil {
		err = utils.CustomErrorWithTrace(err, "reference_id must be unique", 400)
		utils.PanicIfError(err)
	}

	if balanceUpdated > 0 {
		go updateBalance(s.repo, balanceUpdated, userID)
	}

	return dto.WithdrawalResp{
		Withdrawal: dto.Withdrawal{
			ID:          transaction.ID,
			DepositBy:   userID,
			Status:      transaction.Status,
			WithdrawnAt: transaction.TransactionAt,
			Amount:      int64(transaction.Balance),
			ReferenceID: transaction.ReferenceID,
		},
	}
}
