package event

// Event represents the payload provided to the app
type Event struct {
	ChannelNames []string `json:"channel_names"`
	NotifyURL    string   `json:"notify_url"`
}
