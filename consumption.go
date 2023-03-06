package appstore

type ConsumptionRequest struct {
	AccountTenure            int32  `json:"accountTenure"`
	AppAccountToken          string `json:"appAccountToken"`
	ConsumptionStatus        int32  `json:"consumptionStatus"`
	CustomerConsented        bool   `json:"customerConsented"`
	DeliveryStatus           int32  `json:"deliveryStatus"`
	LifetimeDollarsPurchased int32  `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded  int32  `json:"lifetimeDollarsRefunded"`
	Platform                 int32  `json:"platform"`
	PlayTime                 int32  `json:"playTime"`
	SampleContentProvided    bool   `json:"sampleContentProvided"`
	UserStatus               int32  `json:"userStatus"`
}
