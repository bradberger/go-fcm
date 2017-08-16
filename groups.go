package fcm

import "golang.org/x/net/context"

// NotificationGroupResponse is a notification group response. See https://firebase.google.com/docs/cloud-messaging/js/device-group
type NotificationGroupResponse struct {
	NotificationKey string `json:"notification_key"`
}

// NotificationGroupRequest is a notification group request payload. See https://firebase.google.com/docs/cloud-messaging/js/device-group
type NotificationGroupRequest struct {
	Operation           string   `json:"operation"`
	NotificationKeyName string   `json:"notification_key_name"`
	RegistrationIDs     []string `json:"registration_ids"`
}

// NotificationGroup sends a notification group request to the API
func (c *Client) NotificationGroup(ctx context.Context, req *NotificationGroupRequest) (*NotificationGroupResponse, error) {
	var result NotificationGroupResponse
	_, err := c.Post(ctx, "https://android.googleapis.com/gcm/notification", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
