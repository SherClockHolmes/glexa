// JSON interface for alexa Request schema
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/alexa-skills-kit-interface-reference
package glexa

// Request types
const (
	LaunchRequest       = "LaunchRequest"
	IntentRequest       = "IntentRequest"
	SessionEndedRequest = "SessionEndedRequest"
)

type Attributes map[string]struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type Slots map[string]struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type RequestBody struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Request Request `json:"request"`
}

func (r *RequestBody) GetType() string {
	return r.Request.Type
}

type Request struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	Timestamp string `json:"timestamp"`
	Intent    struct {
		Name  string `json:"name"`
		Slots Slots  `json:"slots"`
	} `json:"intent"`
	Reason string `json:"reason"`
}

type Session struct {
	New         bool       `json:"new"`
	SessionID   string     `json:"sessionId"`
	Attributes  Attributes `json:"attributes"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken"`
	} `json:"user"`
}
