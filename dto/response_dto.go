package dto

import "time"

type InitializeAccountResp struct {
	Token string `json:"token"`
}

type Wallet struct {
	ID        string    `json:"id"`
	OwnedBy   string    `json:"owned_by"`
	Status    string    `json:"status"`
	EnabledAt time.Time `json:"enabled_at"`
	Balance   int64     `json:"balance"`
}

type WalletResp struct {
	Wallet Wallet `json:"wallet"`
}

type WalletDisabled struct {
	ID         string    `json:"id"`
	OwnedBy    string    `json:"owned_by"`
	Status     string    `json:"status"`
	DisabledAt time.Time `json:"disabled_at"`
	Balance    int64     `json:"balance"`
}

type WalletDisabledResp struct {
	Wallet WalletDisabled `json:"wallet"`
}

type Transaction struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	TransactedAt time.Time `json:"transacted_at" copier:"TransactionAt"`
	Type         string    `json:"type"`
	Amount       int64     `json:"amount" copier:"Balance"`
	ReferenceID  string    `json:"reference_id"`
}

type TransactionResp struct {
	Transaction []Transaction `json:"transactions"`
}

type Deposit struct {
	ID          string    `json:"id"`
	DepositBy   string    `json:"deposited_by"`
	Status      string    `json:"status"`
	DepositedAt time.Time `json:"deposited_at"`
	Amount      int64     `json:"amount" copier:"Balance"`
	ReferenceID string    `json:"reference_id"`
}

type DepositResp struct {
	Deposit Deposit `json:"deposit"`
}

type Withdrawal struct {
	ID          string    `json:"id"`
	DepositBy   string    `json:"deposited_by"`
	Status      string    `json:"status"`
	WithdrawnAt time.Time `json:"withdrawn_at"`
	Amount      int64     `json:"amount" copier:"Balance"`
	ReferenceID string    `json:"reference_id"`
}

type WithdrawalResp struct {
	Withdrawal Withdrawal `json:"withdrawal"`
}
