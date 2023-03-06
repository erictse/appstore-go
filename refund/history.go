package refund

import (
	"net/url"
)

type History url.Values

type HistoryOption func(*url.Values)

func WithNextToken(revision string) HistoryOption {
	return func(query *url.Values) {
		(*query).Set("revision", revision)
	}
}
