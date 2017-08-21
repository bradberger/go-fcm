package fcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	// MaxTTL the default ttl for a notification
	MaxTTL = 2419200
	// PriorityHigh notification priority
	PriorityHigh = "high"
	// PriorityNormal notification priority
	PriorityNormal = "normal"

	retryAfterHeader = "Retry-After"
	errorKey         = "error"
)

var (
	// HTTPClient can be used to use a custom HTTP client for requests to the API.
	// If nil, then the default http.Client will be used
	HTTPClient *http.Client
	// Timeout is the default timeout for HTTP clients
	Timeout = 10 * time.Second
)

var (
	// retreyableErrors whether the error is a retryable
	retreyableErrors = map[string]bool{
		"Unavailable":         true,
		"InternalServerError": true,
	}

	fcmServerURL = "https://fcm.googleapis.com/fcm/send"
)

// Client stores the key and the Message (Msg)
type Client struct {
	APIKey string
}

type Error struct {
	Error string `json:"error"`
}

// NotificationPayload notification message payload
type NotificationPayload struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	Icon         string `json:"icon,omitempty"`
	Sound        string `json:"sound,omitempty"`
	Badge        string `json:"badge,omitempty"`
	Tag          string `json:"tag,omitempty"`
	Color        string `json:"color,omitempty"`
	ClickAction  string `json:"click_action,omitempty"`
	BodyLocKey   string `json:"body_loc_key,omitempty"`
	BodyLocArgs  string `json:"body_loc_args,omitempty"`
	TitleLocKey  string `json:"title_loc_key,omitempty"`
	TitleLocArgs string `json:"title_loc_args,omitempty"`
}

// NewClient init and create fcm client
func NewClient(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

// apiKeyHeader generates the value of the Authorization key
func (c *Client) apiKeyHeader() string {
	return fmt.Sprintf("key=%v", c.APIKey)
}

// Get sends a HTTP GET request to the urlStr and decodes the response into `out`
func (c *Client) Get(ctx context.Context, urlStr string, out interface{}) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, out)
}

// Do sends an HTTP request and decodes the results into `out`
func (c *Client) Do(ctx context.Context, req *http.Request, out interface{}) (*http.Response, error) {
	if req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", c.apiKeyHeader())
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.getHTTPClient(ctx).Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	// TODO add an error type here.
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return resp, fmt.Errorf("error reading response body: %v", err)
	}
	log.Debugf(ctx, "Body: %s", string(bodyBytes))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	if resp.StatusCode >= http.StatusBadRequest {
		var errResp Error
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil {
			return resp, errors.New(errResp.Error)
		}
		return resp, errors.New(string(bodyBytes))
	}
	if out != nil {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// Post sends a POST request with JSON payload to the urlStr, and decodes the response into out
func (c *Client) Post(ctx context.Context, urlStr string, data interface{}, out interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	if data != nil {
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest("POST", urlStr, &buf)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, out)
}
