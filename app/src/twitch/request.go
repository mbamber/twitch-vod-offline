package twitch

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const twitchAPI = "https://api.twitch.tv/helix"

// TwitchRequest sends a request using a connection
func TwitchRequest(conn *Connection, requestType, path string) (body []byte, err error) {

	// Create a basic request
	req, err := http.NewRequest(requestType, twitchAPI+path, nil)
	if err != nil {
		return
	}

	// Set the client ID if a connection was provided
	if conn != nil {
		req.Header.Set("Client-ID", conn.ClientID)
	}

	// Perform the request
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return
	}

	// Check the return status code is valid
	switch resp.StatusCode {
	case 200:
		break
	default:
		return nil, fmt.Errorf("Invalid reponse from Twitch (Error code: %d)", resp.StatusCode)
	}

	// Close the body when we're done with it
	defer resp.Body.Close()

	// Get the body contents
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return
}
