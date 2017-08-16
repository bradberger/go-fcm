package fcm

import (
	"time"
)

// ResponseStatus represents fcm response message - (tokens and topics)
type ResponseStatus struct {
	Ok           bool
	StatusCode   int
	MulticastID  int64               `json:"multicast_id"`
	Success      int                 `json:"success"`
	Fail         int                 `json:"failure"`
	CanonicalIDs int                 `json:"canonical_ids"`
	Results      []map[string]string `json:"results,omitempty"`
	MsgID        int64               `json:"message_id,omitempty"`
	Err          string              `json:"error,omitempty"`
	RetryAfter   string
}

// IsTimeout check whether the response timeout based on http response status
// code and if any error is retryable
func (r *ResponseStatus) IsTimeout() bool {
	if r.StatusCode >= 500 {
		return true
	}
	if r.StatusCode == 200 {
		for _, val := range r.Results {
			for k, v := range val {
				if k == errorKey && retreyableErrors[v] == true {
					return true
				}
			}
		}
	}

	return false
}

// GetRetryAfterTime converts the retrey after response header
// to a time.Duration
func (r *ResponseStatus) GetRetryAfterTime() (t time.Duration, e error) {
	t, e = time.ParseDuration(r.RetryAfter)
	return
}
