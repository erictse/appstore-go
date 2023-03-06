# appstore-go

A Go client for the App Store Server API

[![Go Report Card](https://goreportcard.com/badge/github.com/erictse/appstore-go)](https://goreportcard.com/report/github.com/erictse/appstore-go)

The [App Store Server API](https://developer.apple.com/documentation/appstoreserverapi) returns
data signed using the JSON Web Signature (JWS) spec. This client implementation abstracts away
JWS for a pleasant developer experience.

## Testing

There aren't automated tests included in this repo because I haven't determined the proper way
to do it, but I'm open to hearing how to remedy that.

I have tested the client using production data for these operations:

- Look Up Order ID
- Get All Subscription Statuses
- Get Transaction History
- Get Notification History
- Get Refund History
- Request a Test Notification (did not verify that a test notification came)
- Get Test Notification Status

I did not test these operations, although they client is able to communicate with the API service
fine:

- Extend a Subscription Renewal Date
- Extend Subscription Renewal Dates for All Active Subscribers
- Get Status of Subscription Renewal Date Extensions

## What’s next

1. Reuse JSON Web Token (JWT) until expired. This client generates a new JWT for each HTTP request sent, however
according to this tip in the
[App Store Connect API docs](https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests#3031059):
_"You do not need to generate a new token for every API request. To get better performance from
the App Store Connect API, reuse the same signed token for multiple requests until it expires."_
1. Support [App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications)
1. Support [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi/)

## Contributing

You're welcome to share any bug reports, questions, and feedback (including pull requests).
