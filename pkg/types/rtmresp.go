package types

// ResponseRtmStart represents rtm start message
type ResponseRtmStart struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	URL   string `json:"url"`
	Bot   BotID  `json:"self"`
}

// BotID represents the object storing the user ID
type BotID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Message represents messages written to and read from web socket.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

//Represents the response from the AWS me
type AccNumRespMsg struct {
	D Data `json:"data"`
}

type Data struct {
	AccountNumber string `json:"AccountNumber"`
	OwnerTeamId   string `json:"OwnerTeamId"`
}
