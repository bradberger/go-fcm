package fcm

import "golang.org/x/net/context"

// APNTokenImportRequest is an APN token import request data structure
type APNTokenImportRequest struct {
	Application string   `json:"application"`
	Sandbox     bool     `json:"sandbox"`
	APNSTokens  []string `json:"apns_tokens"`
}

// APNTokenImportResult is the API response for APNTokenImportRequest
type APNTokenImportResult struct {
	Results []APNResult `json:"results"`
}

// APNResult is a sub-struct of the API response for APNTokenImportRequest
type APNResult struct {
	APNSToken         string `json:"apns_token"`
	Status            string `json:"status"`
	RegistrationToken string `json:"registration_token"`
}

// ImportAPNToken sends an APNTokenImportRequest to the API
func (c *Client) ImportAPNToken(ctx context.Context, req *APNTokenImportRequest) (*APNTokenImportResult, error) {
	var result APNTokenImportResult
	if _, err := c.Post(ctx, "https://iid.googleapis.com/iid/v1:batchImport", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
