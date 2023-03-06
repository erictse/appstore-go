package appstore

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/erictse/appstore-go/notification"
	"github.com/erictse/appstore-go/refund"
	"github.com/erictse/appstore-go/transaction"
	"github.com/golang-jwt/jwt/v4"
)

const (
	hostProd    = "https://api.storekit.itunes.apple.com"
	hostSandbox = "https://api.storekit-sandbox.itunes.apple.com"

	pathNotificationHistory     = "/inApps/v1/notifications/history"
	pathOrderLookup             = "/inApps/v1/lookup/"
	pathRefundHistory           = "/inApps/v2/refund/lookup/"
	pathRequestTestNotification = "/inApps/v1/notifications/test"
	pathSendConsumptionInfo     = "/inApps/v1/transactions/consumption/"
	pathSubscriptionExtend      = "/inApps/v1/subscriptions/extend/"
	pathSubscriptionMassExtend  = "/inApps/v1/subscriptions/extend/mass/"
	pathSubscriptionStatuses    = "/inApps/v1/subscriptions/"
	pathTestNotificationStatus  = "/inApps/v1/notifications/test/"
	pathTransactionHistory      = "/inApps/v1/history/"

	headerAuthorization = "Authorization"
	headerContentType   = "Content-Type"

	contentTypeJSON           = "application/json"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type AppleAPIClaims struct {
	jwt.RegisteredClaims
	BundleID string `json:"bid,omitempty"`
}

type Client struct {
	httpClient *http.Client

	certAppleInterm *x509.Certificate
	certAppleRoot   *x509.Certificate
	claims          *AppleAPIClaims
	host            *string
	keyFunc         jwt.Keyfunc
	keyID           *string
	privateKey      *ecdsa.PrivateKey
	teamID          *string
	verifyOptions   *x509.VerifyOptions
}

type JWSDecoder interface {
	DecodeJWS(jwt.Keyfunc, []byte) error
}

type ClientOption func(*Client)

func WithSandbox() ClientOption {
	return func(c *Client) {
		host := hostSandbox
		c.host = &host
	}
}

func WithClaimsAndKey(bundleID, issuerID, keyID, teamID string, keyPEM []byte) (ClientOption, error) {
	key, parseErr := jwt.ParseECPrivateKeyFromPEM(keyPEM)
	if parseErr != nil {
		return nil, fmt.Errorf("appleapi: could not read private key: %v", parseErr)
	}
	return func(c *Client) {
		c.claims.Issuer = issuerID
		c.claims.BundleID = bundleID
		c.keyID = &keyID
		c.privateKey = key
		c.teamID = &teamID
	}, nil
}

func WithAppleCerts(intermCertPath, rootCertPath string) (ClientOption, error) {
	intermCertDER, err := os.ReadFile(intermCertPath)
	if err != nil {
		return nil, err
	}
	intermCert, err := x509.ParseCertificate(intermCertDER)
	if err != nil {
		return nil, err
	}
	intermPool := x509.NewCertPool()
	intermPool.AddCert(intermCert)

	certAppleRootDER, err := os.ReadFile(rootCertPath)
	if err != nil {
		return nil, err
	}
	rootCert, err := x509.ParseCertificate(certAppleRootDER)
	if err != nil {
		return nil, err
	}
	rootPool := x509.NewCertPool()
	rootPool.AddCert(rootCert)

	return func(c *Client) {
		c.verifyOptions = &x509.VerifyOptions{
			Intermediates: intermPool,
			Roots:         rootPool,
		}
		c.certAppleInterm = intermCert
		c.certAppleRoot = rootCert
	}, nil
}

func NewClient(opts ...ClientOption) (*Client, error) {
	host := hostProd
	keyFunc := func(t *jwt.Token) (any, error) {
		if multi, ok := t.Header["x5c"].([]any); !ok {
			return nil, fmt.Errorf("cert not found in JWS header")
		} else if encoded, ok := multi[0].(string); !ok {
			return nil, fmt.Errorf("invalid cert format in JWS header")
		} else if decoded, err := base64.StdEncoding.DecodeString(encoded); err != nil {
			return nil, fmt.Errorf("unable base64 decode JWS cert: %v", err)
		} else if cert, err := x509.ParseCertificate(decoded); err != nil {
			return nil, fmt.Errorf("could not parse JWS cert: %v", err)
		} else {
			return cert.PublicKey, nil
		}
	}
	c := &Client{
		httpClient: &http.Client{},
		claims:     &AppleAPIClaims{},
		host:       &host,
		keyFunc:    keyFunc,
	}
	c.claims.RegisteredClaims.Audience = []string{"appstoreconnect-v1"}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

func (c Client) endpoint(path string) string {
	return *c.host + path
}

func (c Client) newClaims(issuedAt time.Time) jwt.Claims {
	claims := c.claims
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(issuedAt)
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(issuedAt.Add(time.Minute * 15))
	return claims
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, c.newClaims(time.Now()))
	token.Header["kid"] = c.keyID
	jwtString, err := token.SignedString(c.privateKey)
	if err != nil {
		return nil, err
	}
	header := make(http.Header)
	header.Set(headerContentType, contentTypeFormURLEncoded)
	header.Set(headerAuthorization, fmt.Sprintf("Bearer %s", jwtString))

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	return req, nil
}

func (c *Client) send(req *http.Request, decoder JWSDecoder) error {
	resp, doErr := c.httpClient.Do(req)
	if doErr != nil {
		return doErr
	}
	defer resp.Body.Close()
	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if err := decoder.DecodeJWS(c.keyFunc, data); err != nil {
			return err
		}
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("JWT authorization header is invalid: %s", req.Header.Get(headerAuthorization))
	}
	var payload Error
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("status: %s err: %v", resp.Status, err)
	}
	return payload
}

func (c *Client) ExtendRenewalDate(ctx context.Context, originalTransactionID string,
	body ExtendRenewalDateRequest) (ExtendRenewalDateResponse, error) {

	var r ExtendRenewalDateResponse
	data, err := json.Marshal(body)
	if err != nil {
		return r, err
	}
	uri := c.endpoint(pathSubscriptionExtend + originalTransactionID)
	req, err := c.newRequest(ctx, http.MethodPut, uri, bytes.NewReader(data))
	if err != nil {
		return r, fmt.Errorf("appleapi: client ExtendRenewalDate: %v", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client ExtendRenewalDate: %v", err)
	}
	return r, nil
}

func (c *Client) GetMassExtendRenewalDateStatus(ctx context.Context, productID, requestID string) (MassExtendRenewalDateStatusResponse, error) {
	var r MassExtendRenewalDateStatusResponse
	uri := c.endpoint(pathSubscriptionMassExtend + requestID + "/" + productID)
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client GetMassExtendRenewalDateStatus: %v", err)
	}
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client GetMassExtendRenewalDateStatus: %v", err)
	}
	return r, nil
}

func (c *Client) GetNotificationHistory(ctx context.Context, start, end time.Time, opts ...notification.HistoryOption) (NotificationHistoryResponse, error) {
	var r NotificationHistoryResponse
	query := url.Values{}
	body := notification.HistoryBody{
		StartDate: start.UnixMilli(),
		EndDate:   end.UnixMilli(),
	}
	for _, opt := range opts {
		opt(&query, &body)
	}
	data, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		return r, jsonErr
	}
	buf := bytes.NewBuffer(data)
	uri := c.endpoint(pathNotificationHistory)
	if len(query) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, query.Encode())
	}
	req, err := c.newRequest(ctx, http.MethodPost, uri, buf)
	if err != nil {
		return r, fmt.Errorf("appleapi: client GetNotificationHistory: %v", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client GetNotificationHistory: %v", err)
	}
	return r, nil
}

func (c *Client) GetRefundHistory(ctx context.Context, originalTransactionID string, opts ...refund.HistoryOption) (RefundHistoryResponse, error) {
	var r RefundHistoryResponse
	query := url.Values{}
	for _, opt := range opts {
		opt(&query)
	}
	uri := c.endpoint(pathRefundHistory + originalTransactionID)
	if len(query) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, query.Encode())
	}
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client GetRefundHistory: %v", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client GetRefundHistory: %v", err)
	}
	return r, nil
}

func (c *Client) GetSubscriptionStatuses(ctx context.Context, originalTransactionID string) (StatusResponse, error) {
	var r StatusResponse
	uri := c.endpoint(pathSubscriptionStatuses + originalTransactionID)
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client GetSubscriptionStatuses: %v", err)
	}
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client GetSubscriptionStatuses: %v", err)
	}
	return r, nil
}

func (c *Client) GetTestNotificationStatus(ctx context.Context, token string) (CheckTestNotificationResponse, error) {
	var resp CheckTestNotificationResponse
	uri := c.endpoint(pathTestNotificationStatus + token)
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return resp, fmt.Errorf("appleapi: client RequestTestNotification: %v", err)
	}
	if err := c.send(req, &resp); err != nil {
		return resp, fmt.Errorf("appleapi: client RequestTestNotification: %v", err)
	}
	return resp, nil
}

func (c *Client) GetTransactionHistory(ctx context.Context, originalTransactionID string, opts ...transaction.HistoryOption) (HistoryResponse, error) {
	var r HistoryResponse
	query := url.Values{}
	for _, opt := range opts {
		opt(&query)
	}
	uri := c.endpoint(pathTransactionHistory + originalTransactionID)
	if len(query) > 0 {
		uri = fmt.Sprintf("%s?%s", uri, query.Encode())
	}
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client GetTransactionHistory: %v", err)
	}
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client GetTransactionHistory: %v", err)
	}
	return r, nil
}

func (c *Client) LookupOrder(ctx context.Context, orderID string) (OrderLookupResponse, error) {
	var r OrderLookupResponse
	uri := c.endpoint(pathOrderLookup + orderID)
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client create request LookupOrder: %v", err)
	}
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client send LookupOrder: %v", err)
	}
	return r, nil
}

func (c *Client) MassExtendRenewalDates(ctx context.Context, body MassExtendRenewalDateRequest) (MassExtendRenewalDateResponse, error) {
	var r MassExtendRenewalDateResponse
	data, err := json.Marshal(body)
	if err != nil {
		return r, err
	}
	uri := c.endpoint(pathSubscriptionMassExtend)
	req, err := c.newRequest(ctx, http.MethodPost, uri, bytes.NewReader(data))
	if err != nil {
		return r, fmt.Errorf("appleapi: client MassExtendRenewalDates: %v", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client MassExtendRenewalDates: %v", err)
	}
	return r, nil
}

func (c *Client) RequestTestNotification(ctx context.Context) (SendTestNotificationResponse, error) {
	var r SendTestNotificationResponse
	uri := c.endpoint(pathRequestTestNotification)
	req, err := c.newRequest(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return r, fmt.Errorf("appleapi: client RequestTestNotification: %v", err)
	}
	if err := c.send(req, &r); err != nil {
		return r, fmt.Errorf("appleapi: client RequestTestNotification: %v", err)
	}
	return r, nil
}

func (c *Client) SendConsumptionInfo(ctx context.Context, originalTransactionID string, body ConsumptionRequest) error {
	var r SendTestNotificationResponse
	uri := c.endpoint(pathSendConsumptionInfo + originalTransactionID)
	data, jsonErr := json.Marshal(body)
	if jsonErr != nil {
		return jsonErr
	}
	req, err := c.newRequest(ctx, http.MethodPut, uri, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("appleapi: client SendConsumptionInfo: %v", err)
	}
	req.Header.Set(headerContentType, contentTypeJSON)
	if err := c.send(req, &r); err != nil {
		return fmt.Errorf("appleapi: client SendConsumptionInfo: %v", err)
	}
	return nil
}
