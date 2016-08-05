package glexa

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
	OutputSpeech     *OutputSpeech `json:"outputSpeech"`
	Card             *Card         `json:"card"`
	Reprompt         *Reprompt     `json:"reprompt"`
	ShouldEndSession bool          `json:"shouldEndSession"`
}

// Convenience method for reprompting the user
func (r *Response) Ask(text string) *Response {
	r.Reprompt = &Reprompt{
		OutputSpeech: &OutputSpeech{
			Type: "PlainText",
			Text, text,
		},
	}

	return r
}

// Convenience method for setting the card
func (r *Response) SetCard(card *Card) *Response {
	r.Card = card
	return r
}

// Convenience method for setting the plain text output speech type
func (r *Response) Tell(text string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}

	return r
}

// Convenience method for setting the SSML output speech type
func (r *Response) TellSSML(ssml string) *Response {
	r.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: ssml,
	}

	return r
}
