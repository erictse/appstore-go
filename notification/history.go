package notification

import "net/url"

type HistoryBody struct {
	// Required
	StartDate int64 `json:"startDate"`
	EndDate   int64 `json:"endDate"`

	// Optional
	OriginalTransactionID string `json:"originalTransactionId,omitempty"`
	NotificationSubtype   string `json:"notificationSubtype,omitempty"`
	NotificationType      string `json:"notificationType,omitempty"`
}

type HistoryOption func(*url.Values, *HistoryBody)

func WithOriginalTransactionID(id string) HistoryOption {
	return func(q *url.Values, b *HistoryBody) {
		b.OriginalTransactionID = id
	}
}

func WithType(notificationType string) HistoryOption {
	return func(q *url.Values, b *HistoryBody) {
		b.NotificationType = notificationType
	}
}

func WithSubtype(notificationSubtype string) HistoryOption {
	return func(q *url.Values, b *HistoryBody) {
		b.NotificationSubtype = notificationSubtype
	}
}

func WithNextToken(paginationToken string) HistoryOption {
	return func(q *url.Values, b *HistoryBody) {
		url.Values(*q).Set("paginationToken", paginationToken)
	}
}
