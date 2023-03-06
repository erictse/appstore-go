package appstore

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

type ExtendRenewalDateRequest struct {
	ExtendByDays      int32  `json:"extendByDays"`
	ExtendReasonCode  int32  `json:"extendReasonCode"`
	RequestIdentifier string `json:"requestIdentifier"`
}

type ExtendRenewalDateResponse struct {
	EffectiveDate         Millistamp `json:"effectiveDate"`
	OriginalTransactionId string     `json:"originalTransactionId"`
	Success               bool       `json:"success"`
	WebOrderLineItemId    string     `json:"webOrderLineItemId"`
}

func (p *ExtendRenewalDateResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	return nil
}

type MassExtendRenewalDateRequest struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	ExtendByDays           int32    `json:"extendByDays"`
	ExtendReasonCode       int32    `json:"extendReasonCode"`
	ProductId              string   `json:"productId"`
	StorefrontCountryCodes []string `json:"storefrontCountryCodes"`
}

type MassExtendRenewalDateResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
}

func (p *MassExtendRenewalDateResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	return nil
}

type MassExtendRenewalDateStatusResponse struct {
	RequestIdentifier string     `json:"requestIdentifier"`
	Complete          bool       `json:"complete"`
	CompleteDate      Millistamp `json:"completeDate"`
	FailedCount       int64      `json:"failedCount"`
	SucceededCount    int64      `json:"succeededCount"`
}

func (p *MassExtendRenewalDateStatusResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	return nil
}
