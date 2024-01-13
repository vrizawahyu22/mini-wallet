package dto

type DepositReq struct {
	Amount      float64 `json:"amount"`
	ReferenceID string  `json:"reference_id"`
}

type WithdrawalReq struct {
	Amount      float64 `json:"amount"`
	ReferenceID string  `json:"reference_id"`
}
