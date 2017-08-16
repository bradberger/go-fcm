// +build appengine

package fcm

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func (c *Client) getHTTPClient(ctx context.Context) *http.Client {
	if HTTPClient != nil {
		return HTTPClient
	}
	ctx, _ = context.WithTimeout(ctx, Timeout)
	return urlfetch.Client(ctx)
}
