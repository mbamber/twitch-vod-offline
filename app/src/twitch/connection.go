package twitch

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xfxdev/xlog"
)

// Connection Represents a connection to twitch
type Connection struct {

	// An authorization bearer token. This may be empty and can be refreshed
	// by calling `RefreshToken()`
	Bearer string
	// The twitch client ID used for the connection
	ClientID string
	// The client secret used to get an access token
	ClientSecret string
}

// Pagination represents part of a response with a pagination cursor
type Pagination struct {
	Cursor string `json:"cursor"`
}

// Channels represents a reponse containing a list of channels
type Channels struct {
	Channels []Channel `json:"data"`
}

// Channel represents details about a twitch channel
type Channel struct {
	ID          string `json:"id"`
	ChannelName string `json:"login"`
}

// VODs represents a response containing a list of VODs
type VODs struct {
	VODs       []VOD      `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// VOD represents details about a particular VOD
type VOD struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	DateCreated   string `json:"created_at"`
	DatePublished string `json:"published_at"`
	URL           string `json:"url"`
}

// NewConnection creates a new Connection with a Client ID and Client Secret
func NewConnection(clientID, clientSecret string) Connection {
	conn := Connection{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	xlog.Debugf("Created new connection with ClientID: %s", conn.ClientID)

	return conn
}

// RefreshToken updates the bearer token for the app
func (conn *Connection) RefreshToken() (*Connection, error) {

	path := fmt.Sprintf("/token?client_id=%s&client_secret=%s&grant_type=client_credentials",
		conn.ClientID,
		conn.ClientSecret)

	resp, err := AuthRequest(conn, path)
	if err != nil {
		return nil, err
	}

	// Update the access token
	conn.Bearer = resp.AccessToken
	xlog.Debugf("Updated connection bearer token")

	return conn, nil
}

// GetChannels gets a list of channels from the channel names.
func (conn *Connection) GetChannels(channelNames []string) (channels []Channel, err error) {
	xlog.Debugf("Getting channels for %d channel names: %+v", len(channelNames), channelNames)

	channels = make([]Channel, 0)
	if len(channelNames) == 0 {
		return channels, errors.New("No channel names provided")
	}

	// Build the request URL
	url := "/users?"
	for _, n := range channelNames {
		url = url + "login=" + n + "&"
	}
	url = url[:len(url)-1]
	xlog.Debugf("Generated URL as: %s", url)

	// Perform the request
	resp, err := APIRequest(conn, "GET", url)
	if err != nil {
		return
	}

	jsonChannelIDs := Channels{}
	err = json.Unmarshal(resp.Body, &jsonChannelIDs)
	if err != nil {
		xlog.Errorf("Could not unmarhsal Channel IDs response: %s", resp.Body)
		return
	}
	channels = jsonChannelIDs.Channels
	xlog.Debugf("Successfully unmarshalled %d channels to: %+v", len(channels), channels)

	return
}

// GetVODs gets a list of all VODs created by the given channelIDs.
func (conn *Connection) GetVODs(channelIDs []string) (vods []VOD, err error) {
	xlog.Debugf("Getting VODs for %d channel IDs: %+v", len(channelIDs), channelIDs)

	vods = make([]VOD, 0)
	if len(channelIDs) == 0 {
		return vods, errors.New("No channel IDs provided")
	}

	// Loop through each channel ID and get the VODs for it as we cannot group them all into
	// one API request to twitch
	for _, channelID := range channelIDs {
		url := "/videos?first=100&user_id=" + channelID
		xlog.Debugf("Generated URL as: %s", url)

		// Make the API request
		resp, err := APIRequest(conn, "GET", url)
		if err != nil {
			return nil, err
		}

		jsonVODs := VODs{}
		err = json.Unmarshal(resp.Body, &jsonVODs)
		if err != nil {
			xlog.Errorf("Could not unmarhsal VODs response: %s", resp.Body)
			return vods, err
		}
		xlog.Debugf("Successfully unmarshalled %d VODs to: %+v", len(jsonVODs.VODs), jsonVODs.VODs)

		// Add each VOD to the slice
		for _, v := range jsonVODs.VODs {
			vods = append(vods, v)
		}

		// Handle paginated responses
		for jsonVODs.Pagination.Cursor != "" {
			xlog.Debugf("Making another request to handle pagination. Pagination cursor: %s", jsonVODs.Pagination.Cursor)

			url = "/videos?first=100&after=" + jsonVODs.Pagination.Cursor + "&user_id=" + channelID
			xlog.Debugf("Generated URL as: %s", url)

			// Make the API request
			resp, err := APIRequest(conn, "GET", url)
			if err != nil {
				return nil, err
			}

			jsonVODs = VODs{}
			err = json.Unmarshal(resp.Body, &jsonVODs)
			if err != nil {
				xlog.Errorf("Could not unmarhsal paginated VODs response: %s", resp.Body)
				return vods, err
			}
			xlog.Debugf("Successfully unmarshalled %d paginated VODs to: %+v", len(jsonVODs.VODs), jsonVODs.VODs)

			// Add each VOD to the slice
			for _, v := range jsonVODs.VODs {
				vods = append(vods, v)
			}
		}
	}

	xlog.Debugf("Got %d VODs: %+v", len(vods), vods)
	return vods, nil
}
