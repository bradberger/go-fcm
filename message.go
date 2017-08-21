package fcm

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/appengine/log"
)

// Message represents fcm request message
type Message struct {
	Notification Notification `json:"notification"`
	Priority     string       `json:"priority,omitempty"`
	To           string       `json:"to"`
}

type Notification struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	Icon        string `json:"icon"`
	ClickAction string `json:"click_action"`
}

// TopicResponse is the API response when sending a message to a topic
type TopicResponse struct {
	MessageID    int64                 `json:"message_id"`
	MulticastID  int64                 `json:"multicast_id,omitempty"`
	Success      int                   `json:"success,omitempty"`
	Failure      int                   `json:"failure,omitempty"`
	CanonicalIDs int                   `json:"canonical_ids,omitempty"`
	Results      []TopicResponseResult `json:"results,omitempty"`
}

type TopicResponseResult struct {
	Error string `json:"error"`
}

// Err returns the error message, if any
func (t TopicResponse) Error() error {
	for i := range t.Results {
		if t.Results[i].Error != "" {
			return errors.New(t.Results[i].Error)
		}
	}
	return nil
}

// Send sends a message
func (c *Client) Send(ctx context.Context, m *Message) error {
	m.Priority = ""
	_, err := c.Post(ctx, "https://fcm.googleapis.com/fcm/send", &m, nil)
	return err
}

// Subscribe adds subscriber(s) to the topic
func (c *Client) Subscribe(ctx context.Context, topic string, tokens ...string) error {
	var eg errgroup.Group
	for i := range tokens {
		eg.Go(func(token string) func() error {
			return func() error {
				urlStr := fmt.Sprintf("https://iid.googleapis.com/iid/v1/%s/rel/topics/%s", token, topic)
				log.Debugf(ctx, "Subscribing %s to %s", token, topic)
				resp, err := c.Post(ctx, urlStr, nil, nil)
				log.Debugf(ctx, "Respone: %+v", resp)
				return err
			}
		}(tokens[i]))
	}
	return eg.Wait()
}

// SendTopic sends the message to the subscribers of the topic
func (c *Client) SendTopic(ctx context.Context, m *Message) (*TopicResponse, error) {
	var tr TopicResponse
	log.Debugf(ctx, "Sending message: %+v", m)
	_, err := c.Post(ctx, "https://fcm.googleapis.com/fcm/send", m, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, tr.Error()
}
