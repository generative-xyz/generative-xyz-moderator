package request

import "errors"

type WithDrawItemRequest struct {
	Amount         *string `json:"amount"`
	PaymentType    *string `json:"paymentType"` //referal or project
	WithdrawType   *string `json:"type"`
	ID   *string `json:"id"` //referal  (referal: refereeID, project: tokenID)
}

func (w WithDrawItemRequest) SelfValidate() error {
	if w.Amount == nil {
		return errors.New("Amount is required")
	}else{
		if *w.Amount == "" {
			return errors.New("Amount is not empty")
		}
	}
	
	if w.PaymentType == nil {
		return errors.New("PaymentType is required")
	}else{
		if *w.PaymentType == "" {
			return errors.New("PaymentType is not empty")
		}
	}
	
	if w.WithdrawType == nil {
		return errors.New("WithdrawType is required")
	}else{
		if *w.WithdrawType == "" {
			return errors.New("WithdrawType is not empty")
		}
	}
	
	if w.ID == nil {
		return errors.New("Withdraw ID is required")
	}else{
		if *w.ID == "" {
			return errors.New("ID is not empty")
		}
	}
	
	return nil
}