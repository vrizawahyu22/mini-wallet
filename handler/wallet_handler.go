package handler

import (
	"mini-wallet/dto"
	"mini-wallet/service"
	"mini-wallet/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type WalletHandler interface {
	InitializeAccount(w http.ResponseWriter, r *http.Request)
	EnableWallet(w http.ResponseWriter, r *http.Request)
	DisableWallet(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
	GetTransactions(w http.ResponseWriter, r *http.Request)
	DepositWallet(w http.ResponseWriter, r *http.Request)
	WithdrawalWallet(w http.ResponseWriter, r *http.Request)

	SetupWalletRoutes(route *chi.Mux)
}

type WalletHandlerImpl struct {
	walletSvc service.WalletSvc
}

func NewWalletHandler(
	walletSvc service.WalletSvc,
) WalletHandler {
	return &WalletHandlerImpl{
		walletSvc: walletSvc,
	}
}

func (h *WalletHandlerImpl) InitializeAccount(w http.ResponseWriter, r *http.Request) {
	key := "customer_xid"
	customerXID := r.FormValue(key)

	if customerXID == "" {
		utils.PanicValidationError([]utils.ValidationError{
			{
				Key:     key,
				Message: "Missing data for required field.",
			},
		}, 400)
	}

	resp := h.walletSvc.InitializeAccount(r.Context(), customerXID)

	utils.GenerateSuccessResp(w, resp, 201)
}

func (h *WalletHandlerImpl) EnableWallet(w http.ResponseWriter, r *http.Request) {
	resp := h.walletSvc.EnableWallet(r.Context())

	utils.GenerateSuccessResp(w, resp, 201)
}

func (h *WalletHandlerImpl) DisableWallet(w http.ResponseWriter, r *http.Request) {
	key := "is_disabled"
	isDisabled := r.FormValue(key)

	if isDisabled == "" {
		utils.PanicValidationError([]utils.ValidationError{
			{
				Key:     key,
				Message: "Missing data for required field.",
			},
		}, 400)
	}

	isDisabledBool, err := strconv.ParseBool(isDisabled)
	if err != nil {
		utils.PanicValidationError([]utils.ValidationError{
			{
				Key:     key,
				Message: err.Error(),
			},
		}, 400)
	}

	resp := h.walletSvc.DisableWallet(r.Context(), isDisabledBool)

	utils.GenerateSuccessResp(w, resp, 200)
}

func (h *WalletHandlerImpl) GetWallet(w http.ResponseWriter, r *http.Request) {
	resp := h.walletSvc.GetWallet(r.Context())

	utils.GenerateSuccessResp(w, resp, 200)
}

func (h *WalletHandlerImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	resp := h.walletSvc.GetTransactions(r.Context())

	utils.GenerateSuccessResp(w, resp, 200)
}

func (h *WalletHandlerImpl) DepositWallet(w http.ResponseWriter, r *http.Request) {
	keyAmount := "amount"
	keyReferenceID := "reference_id"
	amount := r.FormValue(keyAmount)
	referenceID := r.FormValue(keyReferenceID)

	var validationErrors []utils.ValidationError
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if amount == "" || err != nil {
		validationErrors = append(validationErrors, utils.ValidationError{
			Key:     keyAmount,
			Message: "Missing data for required field.",
		})
	}
	if referenceID == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Key:     keyReferenceID,
			Message: "Missing data for required field.",
		})
	}

	if len(validationErrors) > 0 {
		utils.PanicValidationError(validationErrors, 400)
	}

	resp := h.walletSvc.DepositWallet(r.Context(), dto.DepositReq{
		Amount:      amountFloat,
		ReferenceID: referenceID,
	})

	utils.GenerateSuccessResp(w, resp, 200)
}

func (h *WalletHandlerImpl) WithdrawalWallet(w http.ResponseWriter, r *http.Request) {
	keyAmount := "amount"
	keyReferenceID := "reference_id"
	amount := r.FormValue(keyAmount)
	referenceID := r.FormValue(keyReferenceID)

	var validationErrors []utils.ValidationError
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if amount == "" || err != nil {
		validationErrors = append(validationErrors, utils.ValidationError{
			Key:     keyAmount,
			Message: "Missing data for required field.",
		})
	}
	if referenceID == "" {
		validationErrors = append(validationErrors, utils.ValidationError{
			Key:     keyReferenceID,
			Message: "Missing data for required field.",
		})
	}

	if len(validationErrors) > 0 {
		utils.PanicValidationError(validationErrors, 400)
	}

	resp := h.walletSvc.WithdrawalWallet(r.Context(), dto.WithdrawalReq{
		Amount:      amountFloat,
		ReferenceID: referenceID,
	})

	utils.GenerateSuccessResp(w, resp, 200)
}

func (h *WalletHandlerImpl) SetupWalletRoutes(route *chi.Mux) {
	route.Post("/api/v1/init", h.InitializeAccount)
	route.Post("/api/v1/wallet", utils.CheckIsAuthenticated(h.EnableWallet))
	route.Get("/api/v1/wallet", utils.CheckIsAuthenticated(h.GetWallet))
	route.Patch("/api/v1/wallet", utils.CheckIsAuthenticated(h.DisableWallet))
	route.Get("/api/v1/wallet/transactions", utils.CheckIsAuthenticated(h.GetTransactions))
	route.Post("/api/v1/wallet/deposits", utils.CheckIsAuthenticated(h.DepositWallet))
	route.Post("/api/v1/wallet/withdrawals", utils.CheckIsAuthenticated(h.WithdrawalWallet))
}
