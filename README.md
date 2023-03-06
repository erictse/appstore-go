# appstore-go

A Go client for the App Store Server API

[![Go Report Card](https://goreportcard.com/badge/github.com/erictse/appstore-go)](https://goreportcard.com/report/github.com/erictse/appstore-go)

The [App Store Server API](https://developer.apple.com/documentation/appstoreserverapi) returns data signed using the JSON Web Signature (JWS) spec. This client implementation abstracts away JWS for a pleasant developer experience.

## Quickstart

### Requirements

- Bundle ID - see App Information under General group for the app detail page in App Store Connect.
- Team ID - see [Developer Membership Details](https://developer.apple.com/account#MembershipDetailsCard)
- Get Apple public certificates from the [Apple PKI webpage](https://www.apple.com/certificateauthority/).
- Create and download a private key the [Users and Access Keys page](https://appstoreconnect.apple.com/access/api) in App Store Connect.
- Use the Issuer ID and Key ID from that page as well. While there are two key types, App Store Connect API and In-App Purchase, only use the In-App Purchase key for this client since it does not support the App Store Connect API.

### Configure client

```go
// The private key should be in PEM format beginning with -----BEGIN PRIVATE KEY----- 
// and ending with -----END PRIVATE KEY-----
keyData, keyErr := os.ReadFile("~/.secret/SubscriptionKey_ABC0DE1F23.p8")
if keyErr != nil {
    log.Fatalln("Could not read private key:", keyErr)
}
optCerts, certErr := appstore.WithAppleCerts(
    "~/public/AppleRootCA-G3.cer", // Root certificate
    "~/public/AppleWWDRCAG6.cer",  // Intermediate certificate
)
if certErr != nil {
    log.Fatalln("Could not load Apple certs:", certErr)
}
optClaimsKey, claimsErr := appstore.WithClaimsAndKey(
    "com.example.appname",                  // Bundle ID
    "5bf1bb6b-ddf0-417e-b12e-e9fadd5fc611", // Issuer ID
    "ABC0DE1F23",                           // Key ID
    "1234AB5678",                           // Team ID
    keyData,
)
if claimsErr != nil {
    log.Fatalln("Could not configure client:", claimsErr)
}
client, clientErr := appstore.NewClient(optCerts, optClaimsKey)
if clientErr != nil {
    log.Fatalln("Could not create JWT client:", clientErr)
}
```

## Examples

### Call API with optional parameters

```go
startOpt := transaction.WithStartDate(time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC))
endOpt := transaction.WithEndDate(time.Date(2023, time.March, 5, 0, 0, 0, 0, time.UTC))

resp, err := client.GetTransactionHistory(context.TODO(), "123456789012345", startOpt, endOpt)
if err != nil {
    log.Println("error", err)
}
log.Println(resp.Transactions[0])
```

### Paginate

```go
ctx := context.TODO()
originalTransactionID := "123456789012345"

resp, err := client.GetTransactionHistory(ctx, originalTransactionID)
if err != nil {
    log.Println("error", err)
}

for {
    for _, t := range resp.Transactions {
        log.Println(t)
    }
    if !resp.HasMore {
        break
    }
    nextOpt := transaction.WithNextToken(resp.Revision)
    resp, err = client.GetTransactionHistory(ctx, originalTransactionID, nextOpt)
    if err != nil {
        log.Println("error", err)
    }
}
```

## Testing

There aren't automated tests included in this repo because I haven't determined the proper way to do it, but I'm open to hearing how to remedy that.

I have tested the client using production data for these operations:

- Look Up Order ID
- Get All Subscription Statuses
- Get Transaction History
- Get Notification History
- Get Refund History
- Request a Test Notification (did not verify that a test notification came)
- Get Test Notification Status

I did not test these operations, although they client is able to communicate with the API service fine:

- Extend a Subscription Renewal Date
- Extend Subscription Renewal Dates for All Active Subscribers
- Get Status of Subscription Renewal Date Extensions

## What’s next

1. Reuse JSON Web Token (JWT) until expired. This client generates a new JWT for each HTTP request sent, however according to this tip in the
[App Store Connect API docs](https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests#3031059): _"You do not need to generate a new token for every API request. To get better performance from the App Store Connect API, reuse the same signed token for multiple requests until it expires."_
1. Support [App Store Server Notifications](https://developer.apple.com/documentation/appstoreservernotifications)
1. Support [App Store Connect API](https://developer.apple.com/documentation/appstoreconnectapi/)

## Contributing

You're welcome to share any bug reports, questions, and feedback (including pull requests).
