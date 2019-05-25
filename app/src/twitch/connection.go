package twitch

import (
	"encoding/json"
	"errors"
)

// Represents a connection to twitch
type Connection struct {

	// The twitch client ID used for the connection
	ClientID string
}

// Unmarhsal the body
type ChannelIDs struct {
	Channels []Channel `json:"data"`
}

type Channel struct {
	ID          string `json:"id"`
	ChannelName string `json:"login"`
}

// GetChannels gets a list of channel IDs from the channel names.
// The IDs are returned as a map from name to ID
func (conn *Connection) GetChannels(channelNames []string) (channels []Channel, err error) {

	channels = make([]Channel, 0)
	if len(channelNames) == 0 {
		return channels, errors.New("No channel names provided!")
	}

	// Build the request URL
	url := "/users?"
	for _, n := range channelNames {
		url = url + "login=" + n + "&"
	}
	url = url[:len(url)-1]

	// Perform the request
	body, err := TwitchRequest(conn, "GET", url)
	if err != nil {
		return
	}

	jsonChannelIDs := ChannelIDs{}
	err = json.Unmarshal(body, &jsonChannelIDs)
	if err != nil {
		return
	}
	channels = jsonChannelIDs.Channels

	return
}
