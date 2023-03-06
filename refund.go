package appstore

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

type RefundHistoryResponse struct {
	HasMore            bool      `json:"hasMore"`
	Revision           string    `json:"revision"`
	SignedTransactions []JWSData `json:"signedTransactions"`

	Transactions []JWSTransactionDecodedPayload
}

func (p *RefundHistoryResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	for i, s := range p.SignedTransactions {
		if err := s.Decode(keyFunc, &p.Transactions[i]); err != nil {
			return err
		}
	}
	return nil
}
