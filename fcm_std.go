// +build !appengine

package fcm

import (
	"net/http"

	"golang.org/x/net/context"
)

func (c *Client) getHTTPClient(ctx context.Context) *http.Client {
	if HTTPClient != nil {
		return HTTPClient
	}
	return &http.Client{Timeout: Timeout}
}
