package appstore

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type OrderLookupResponse struct {
	Status             int32     `json:"status"`
	SignedTransactions []JWSData `json:"signedTransactions"`

	Transactions []JWSTransactionDecodedPayload
}

func (r *OrderLookupResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	var errs []error
	for _, t := range r.SignedTransactions {
		var transaction JWSTransactionDecodedPayload
		if err := t.Decode(keyFunc, &transaction); err != nil {
			errs = append(errs, err)
		}
		r.Transactions = append(r.Transactions, transaction)
	}
	if len(errs) > 0 {
		return fmt.Errorf("could not decode all transactions: %v", errors.Join(errs...))
	}
	return nil
}
