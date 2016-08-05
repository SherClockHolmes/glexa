// JSON interface for alexa Response schema
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/alexa-skills-kit-interface-reference
package glexa

const (
	outputSpeechPlainText = "PlainText"
	outputSpeechSSML      = "SSML"
)

type ResponseBody struct {
	Version           string      `json:"version"`
	SessionAttributes *Attributes `json:"sessionAttributes"`
	Response          *Response   `json:"response"`
}

type Card struct {
	Type    string `json:"type"` // Simple, Standard, LinkAccount
	Title   string `json:"title"`
	Content string `json:"content"`
	Text    string `json:"text"`
	Image   *Image `json:"image"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl"`
	LargeImageURL string `json:"largeImageUrl"`
}

type OutputSpeech struct {
	Type string `json:"type"` // PlainText, SSML
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type Reprompt struct {
	OutputSpeech *OutputSpeech `json:"outputSpeech"`
}

type Response struct {
	Error            error         `json:"-"`
	OutputSpeech     *OutputSpeech `json:"outputSpeech"`
	Card             *Card         `json:"card"`
	Reprompt         *Reprompt     `json:"reprompt"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

func NewResponse() *Response {
	return &Response{}
}

// Reprompt the user with text
func (r *Response) Ask(text string) *Response {
	r.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: outputSpeechPlainText,
			Text: text,
		},
	}

	return r
}

// Reprompt the user with SSML
func (r *Response) AskSSML(ssml string) *Response {
	r.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: outputSpeechSSML,
			SSML: ssml,
		},
	}

	return r
}

// Set the card
func (r *Response) SetCard(card *Card) *Response {
	r.Card = card
	return r
}

// Set the plain text output speech
func (r *Response) Tell(text string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: outputSpeechPlainText,
		Text: text,
	}

	return r
}

// Set the SSML output speech
func (r *Response) TellSSML(ssml string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: outputSpeechSSML,
		SSML: ssml,
	}

	return r
}
