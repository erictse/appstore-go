package transaction

import (
	"net/url"
	"strconv"
	"time"
)

type HistoryOption func(*url.Values)

func WithStartDate(t time.Time) HistoryOption {
	return func(query *url.Values) {
		query.Set("startDate", strconv.FormatInt(t.UnixMilli(), 10))
	}
}

func WithEndDate(t time.Time) HistoryOption {
	return func(query *url.Values) {
		query.Set("endDate", strconv.FormatInt(t.UnixMilli(), 10))
	}
}

func WithNextToken(revision string) HistoryOption {
	return func(query *url.Values) {
		query.Set("revision", revision)
	}
}

func WithProductIDs(ids ...string) HistoryOption {
	return func(query *url.Values) {
		params := *query
		for _, i := range ids {
			params.Add("productId", i)
		}
	}
}

func WithProductTypes(types ...string) HistoryOption {
	return func(query *url.Values) {
		for _, t := range types {
			query.Add("productType", t)
		}
	}
}

func WithExcludeRevoked(excludeRevoked bool) HistoryOption {
	return func(query *url.Values) {
		val := "false"
		if excludeRevoked {
			val = "true"
		}
		query.Set("excludeRevoked", val)
	}
}
