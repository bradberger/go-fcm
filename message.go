package fcm

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

// Message represents fcm request message
type Message struct {
	Notification struct {
		Title       string `json:"title"`
		Body        string `json:"body"`
		Icon        string `json:"icon"`
		ClickAction string `json:"click_action"`
	} `json:"notification"`
	Priority string `json:"priority,omitempty"`
	To       string `json:"to"`
}

// TopicResponse is the API response when sending a message to a topic
type TopicResponse struct {
	MessageID string `json:"message_id"`
	Error     string `json:"error"`
}

// Err returns the error message, if any
func (t TopicResponse) Err() error {
	if t.Error == "" {
		return nil
	}
	return errors.New(t.Error)
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
				_, err := c.Post(ctx, fmt.Sprintf("https://iid.googleapis.com/iid/v1/%s/rel/topics/%s", token, topic), nil, nil)
				return err
			}
		}(tokens[i]))
	}
	return eg.Wait()
}

// SendTopic sends the message to the subscribers of the topic
func (c *Client) SendTopic(ctx context.Context, m *Message) (*TopicResponse, error) {
	var tr TopicResponse
	_, err := c.Post(ctx, "https://fcm.googleapis.com/fcm/send", m, &tr)
	if err != nil {
		return nil, err
	}
	return &tr, tr.Err()
}
