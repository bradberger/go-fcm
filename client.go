package fcm

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

const (
	instanceIDInfoWithDetailsSrvURL  = "https://iid.googleapis.com/iid/info/%s?details=true"
	instanceIDInfoNoDetailsSrvURL    = "https://iid.googleapis.com/iid/info/%s"
	subscribeInstanceidToTopicSrvURL = "https://iid.googleapis.com/iid/v1/%s/rel/topics/%s"
	batchAddSrvURL                   = "https://iid.googleapis.com/iid/v1:batchAdd"
	batchRemSrvURL                   = "https://iid.googleapis.com/iid/v1:batchRemove"
	apnsBatchImportSrvURL            = "https://iid.googleapis.com/iid/v1:batchImport"
	apnsTokenKey                     = "apns_token"
	statusKey                        = "status"
	regTokenKey                      = "registration_token"
	topics                           = "/topics/"
)

// InstanceIDInfoResponse response for instance id info request
type InstanceIDInfoResponse struct {
	Application        string                                  `json:"application,omitempty"`
	AuthorizedEntity   string                                  `json:"authorizedEntity,omitempty"`
	ApplicationVersion string                                  `json:"applicationVersion,omitempty"`
	AppSigner          string                                  `json:"appSigner,omitempty"`
	AttestStatus       string                                  `json:"attestStatus,omitempty"`
	Platform           string                                  `json:"platform,omitempty"`
	ConnectionType     string                                  `json:"connectionType,omitempty"`
	ConnectDate        string                                  `json:"connectDate,omitempty"`
	Error              string                                  `json:"error,omitempty"`
	Rel                map[string]map[string]map[string]string `json:"rel,omitempty"`
}

// SubscribeResponse response for single topic subscribtion
type SubscribeResponse struct {
	Error      string `json:"error,omitempty"`
	Status     string
	StatusCode int
}

// BatchRequest add/remove request
type BatchRequest struct {
	To        string   `json:"to,omitempty"`
	RegTokens []string `json:"registration_tokens,omitempty"`
}

// BatchResponse add/remove response
type BatchResponse struct {
	Error      string              `json:"error,omitempty"`
	Results    []map[string]string `json:"results,omitempty"`
	Status     string
	StatusCode int
}

// ApnsBatchRequest apns import request
type ApnsBatchRequest struct {
	App        string   `json:"application,omitempty"`
	Sandbox    bool     `json:"sandbox,omitempty"`
	ApnsTokens []string `json:"apns_tokens,omitempty"`
}

// ApnsBatchResponse apns import response
type ApnsBatchResponse struct {
	Results    []map[string]string `json:"results,omitempty"`
	Error      string              `json:"error,omitempty"`
	Status     string
	StatusCode int
}

// GetInfo gets the instance id info
func (c *Client) GetInfo(ctx context.Context, withDetails bool, instanceIDToken string) (*InstanceIDInfoResponse, error) {

	urlStr := generateGetInfoURL(instanceIDInfoNoDetailsSrvURL, instanceIDToken)
	if withDetails == true {
		urlStr = generateGetInfoURL(instanceIDInfoWithDetailsSrvURL, instanceIDToken)
	}

	var infoResponse InstanceIDInfoResponse
	if _, err := c.Get(ctx, urlStr, &infoResponse); err != nil {
		return nil, err
	}
	return &infoResponse, nil
}

// generateGetInfoUrl generate based on with details and the instance token
func generateGetInfoURL(srv, instanceIDToken string) string {
	return fmt.Sprintf(srv, instanceIDToken)
}

// SubscribeToTopic subscribes a single device/token to a topic
func (c *Client) SubscribeToTopic(ctx context.Context, instanceIDToken, topic string) (*SubscribeResponse, error) {

	var subResp SubscribeResponse
	r, err := c.Post(ctx, generateSubToTopicURL(instanceIDToken, topic), nil, &subResp)
	if err != nil {
		return nil, err
	}
	subResp.Status = r.Status
	subResp.StatusCode = r.StatusCode
	return &subResp, nil
}

// generateSubToTopicURL generates a url based on the instnace id and topic name
func generateSubToTopicURL(instanceID, topic string) string {
	Tmptopic := strings.ToLower(topic)
	if strings.Contains(Tmptopic, "/topics/") {
		tmp := strings.Split(topic, "/")
		topic = tmp[len(tmp)-1]
	}
	return fmt.Sprintf(subscribeInstanceidToTopicSrvURL, instanceID, topic)
}

// BatchSubscribeToTopic subscribes (many) devices/tokens to a given topic
func (c *Client) BatchSubscribeToTopic(ctx context.Context, tokens []string, topic string) (*BatchResponse, error) {
	var result BatchResponse
	r, err := c.Post(ctx, batchAddSrvURL, &BatchRequest{To: topics + extractTopicName(topic), RegTokens: tokens}, &result)
	if err != nil {
		return nil, err
	}
	result.Status = r.Status
	result.StatusCode = r.StatusCode
	return &result, nil
}

// BatchUnsubscribeFromTopic unsubscribes (many) devices/tokens from a given topic
func (c *Client) BatchUnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*BatchResponse, error) {

	var result BatchResponse
	data := &BatchRequest{To: topics + extractTopicName(topic), RegTokens: tokens}
	r, err := c.Post(ctx, batchRemSrvURL, data, result)
	if err != nil {
		return nil, err
	}
	result.Status = r.Status
	result.StatusCode = r.StatusCode
	return &result, nil
}

// extractTopicName extract topic name for valid topic name input
func extractTopicName(inTopic string) (result string) {
	Tmptopic := strings.ToLower(inTopic)
	if strings.Contains(Tmptopic, "/topics/") {
		tmp := strings.Split(inTopic, "/")
		result = tmp[len(tmp)-1]
		return
	}

	result = inTopic
	return
}

// ApnsBatchImportRequest apns import requst
func (c *Client) ApnsBatchImportRequest(ctx context.Context, apnsReq *ApnsBatchRequest) (*ApnsBatchResponse, error) {
	var result ApnsBatchResponse
	r, err := c.Post(ctx, apnsBatchImportSrvURL, apnsReq, &result)
	if err != nil {
		return nil, err
	}
	result.Status = r.Status
	result.StatusCode = r.StatusCode
	return &result, nil
}
