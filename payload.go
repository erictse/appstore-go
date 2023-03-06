package appstore

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Millistamp struct{ time.Time }

func (m *Millistamp) UnmarshalJSON(data []byte) error {
	var val int64
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}
	if val > 0 {
		m.Time = time.UnixMilli(val)
	}
	return nil
}

type JWSData string

func (j JWSData) Decode(keyFunc jwt.Keyfunc, claims jwt.Claims) error {
	if _, err := jwt.ParseWithClaims(string(j), claims, keyFunc); err != nil {
		return err
	}
	return nil
}

type JWSRenewalInfoDecodedPayload struct {
	AutoRenewProductId     string `json:"autoRenewProductId"`
	AutoRenewStatus        int32  `json:"autoRenewStatus"`
	Environment            string `json:"environment"`
	ExpirationIntent       int32  `json:"expirationIntent"`
	IsInBillingRetryPeriod bool   `json:"isInBillingRetryPeriod"`
	OfferIdentifier        string `json:"offerIdentifier"`
	OfferType              int32  `json:"offerType"`
	OriginalTransactionId  string `json:"originalTransactionId"`
	PriceIncreaseStatus    int32  `json:"priceIncreaseStatus"`
	ProductId              string `json:"productId"`

	GracePeriodExpiresDate      *Millistamp `json:"gracePeriodExpiresDate"`
	RecentSubscriptionStartDate *Millistamp `json:"recentSubscriptionStartDate"`
	SignedDate                  *Millistamp `json:"signedDate"`
}

func (p JWSRenewalInfoDecodedPayload) Valid() error {
	return nil
}

type JWSTransactionDecodedPayload struct {
	TransactionID               string `json:"transactionId,omitempty"`
	OriginalTransactionID       string `json:"originalTransactionId,omitempty"`
	WebOrderLineItemID          string `json:"webOrderLineItemId,omitempty"`
	BundleID                    string `json:"bundleId,omitempty"`
	ProductID                   string `json:"productId,omitempty"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier,omitempty"`
	Quantity                    int    `json:"quantity,omitempty"`
	Type                        string `json:"type,omitempty"`
	InAppOwnershipType          string `json:"inAppOwnershipType,omitempty"`
	Environment                 string `json:"environment,omitempty"`

	PurchaseDate         *Millistamp `json:"purchaseDate"`
	OriginalPurchaseDate *Millistamp `json:"originalPurchaseDate"`
	ExpiresDate          *Millistamp `json:"expiresDate"`
	SignedDate           *Millistamp `json:"signedDate"`
}

func (p JWSTransactionDecodedPayload) Valid() error {
	return nil
}
