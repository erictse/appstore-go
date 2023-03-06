package appstore

import (
	"encoding/json"

	"github.com/golang-jwt/jwt/v4"
)

type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}

func (p *SendTestNotificationResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	return nil
}

type CheckTestNotificationResponse struct {
	FirstSendAttemptResult string  `json:"firstSendAttemptResult"`
	SignedPayload          JWSData `json:"signedPayload"`

	Payload ResponseBodyV2DecodedPayload
}

func (p *CheckTestNotificationResponse) DecodeJWS(keyFunc jwt.Keyfunc, data []byte) error {
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}
	if err := p.SignedPayload.Decode(keyFunc, &p.Payload); err != nil {
		return err
	}
	return nil
}
