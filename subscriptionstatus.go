package appstore

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

type LastTransactionsItem struct {
	OriginalTransactionId string  `json:"originalTransactionId"`
	Status                int32   `json:"status"`
	SignedRenewalInfo     JWSData `json:"signedRenewalInfo"`
	SignedTransactionInfo JWSData `json:"signedTransactionInfo"`

	RenewalInfo     JWSRenewalInfoDecodedPayload
	TransactionInfo JWSTransactionDecodedPayload
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                  `json:"subscriptionGroupIdentifier"`
	LastTransactions            []*LastTransactionsItem `json:"lastTransactions"`
}

type StatusResponse struct {
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
	Environment string                            `json:"environment"`
	AppAppleID  int64                             `json:"appAppleId"`
	BundleID    string                            `json:"bundleId"`
}

func (r *StatusResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return err
	}
	for _, group := range r.Data {
		for _, t := range group.LastTransactions {
			if t.SignedRenewalInfo != "" {
				if err := t.SignedRenewalInfo.Decode(keyFunc, &t.RenewalInfo); err != nil {
					return err
				}
			}
			if t.SignedTransactionInfo != "" {
				if err := t.SignedTransactionInfo.Decode(keyFunc, &t.TransactionInfo); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
