package appstore

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

type ResponseBodyV2DecodedPayload struct {
	NotificationType string                               `json:"notificationType"`
	Subtype          string                               `json:"subtype"`
	Data             *ResponseBodyV2DecodedPayloadData    `json:"data,omitempty"`
	Summary          *ResponseBodyV2DecodedPayloadSummary `json:"summary,omitempty"`
	Version          string                               `json:"version"`
	NotificationUUID string                               `json:"notificationUUID"`

	// SignedDate *time.Time
	SignedDate Millistamp `json:"signedDate"`
}

func (p ResponseBodyV2DecodedPayload) Valid() error {
	return nil
}

type ResponseBodyV2DecodedPayloadData struct {
	AppAppleId            int64   `json:"appAppleId"`
	BundleId              string  `json:"bundleId"`
	BundleVersion         string  `json:"bundleVersion"`
	Environment           string  `json:"environment"`
	SignedRenewalInfo     JWSData `json:"signedRenewalInfo"`
	SignedTransactionInfo JWSData `json:"signedTransactionInfo"`

	RenewalInfo     JWSRenewalInfoDecodedPayload
	TransactionInfo JWSTransactionDecodedPayload
}

func (p *ResponseBodyV2DecodedPayloadData) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	if err := p.SignedRenewalInfo.Decode(keyFunc, &p.RenewalInfo); err != nil {
		return err
	}
	if err := p.SignedTransactionInfo.Decode(keyFunc, &p.TransactionInfo); err != nil {
		return err
	}
	return nil
}

type ResponseBodyV2DecodedPayloadSummary struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	Environment            string   `json:"environment"`
	AppAppleId             string   `json:"appAppleId"`
	BundleId               string   `json:"bundleId"`
	ProductId              string   `json:"productId"`
	StorefrontCountryCodes []string `json:"storefrontCountryCodes"`
	FailedCount            int64    `json:"failedCount"`
	SucceededCount         int64    `json:"succeededCount"`
}

type NotificationHistoryResponse struct {
	NotificationHistory []*NotificationHistoryResponseItem `json:"notificationHistory"`
	HasMore             bool                               `json:"hasMore"`
	PaginationToken     string                             `json:"paginationToken"`
}

func (p *NotificationHistoryResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	for _, n := range p.NotificationHistory {
		if err := n.SignedPayload.Decode(keyFunc, &n.Payload); err != nil {
			return err
		}
	}
	return nil
}

type NotificationHistoryResponseItem struct {
	FirstSendAttemptResult string  `json:"firstSendAttemptResult"`
	SignedPayload          JWSData `json:"signedPayload"`

	Payload ResponseBodyV2DecodedPayload
}
