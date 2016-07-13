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

func (r *Response) TellSSML(ssml string) {
	r.OutputSpeech = &OutputSpeech{
		Type: "SSML",
		SSML: ssml,
	}
}

func (r *Response) Tell(text string) {
	r.OutputSpeech = &OutputSpeech{
		Type: "PlainText",
		Text: text,
	}
}
