package appstore

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type HistoryResponse struct {
	AppAppleId         int64     `json:"appAppleId"`
	BundleId           string    `json:"bundleId"`
	Environment        string    `json:"environment"`
	HasMore            bool      `json:"hasMore"`
	Revision           string    `json:"revision"`
	SignedTransactions []JWSData `json:"signedTransactions"`

	Transactions []JWSTransactionDecodedPayload
}

func (r *HistoryResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	var errs []error
	for _, s := range r.SignedTransactions {
		var transaction JWSTransactionDecodedPayload
		if err := s.Decode(keyFunc, &transaction); err != nil {
			errs = append(errs, err)
		}
		r.Transactions = append(r.Transactions, transaction)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
