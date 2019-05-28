package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xfxdev/xlog"
)

// AuthRequest makes a request to the Twitch authentication endpoint
func AuthRequest(conn *Connection, path string) (resp AuthResponse, err error) {
	const twitchAuth = "https://id.twitch.tv/oauth2"
	url := fmt.Sprintf("%s%s", twitchAuth, path)

	body, code, headers, err := Request(conn, "POST", url, nil)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp, err
	}

	resp.Headers = headers
	resp.ResponseCode = code

	return resp, nil
}

// AuthResponse repesents a response returned from an AuthRequest
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	Headers      http.Header `json:"headers"`
	ResponseCode int         `json:"response_code"`
}

// APIRequest makes a request to the Twitch API endpoint
func APIRequest(conn *Connection, requestType, path string) (resp APIResponse, err error) {
	const twitchAPI = "https://api.twitch.tv/helix"
	url := fmt.Sprintf("%s%s", twitchAPI, path)

	// Construct the headers
	headers := make(map[string]string, 0)

	// Check there is a Bearer token set, and refresh it if not
	if conn.Bearer == "" {
		xlog.Debugf("Refreshing Bearer before creating API request")
		_, err := conn.RefreshToken()
		if err != nil {
			return resp, err
		}
	}
	headers["Authorization"] = fmt.Sprintf("Bearer %s", conn.Bearer)

	body, responseCode, responseHeaders, err := Request(conn, requestType, url, headers)
	if err != nil {
		return resp, err
	}

	resp.Body = body
	resp.ResponseCode = responseCode
	resp.Headers = responseHeaders

	return
}

// APIResponse represents a response returned from an APIRequest
type APIResponse struct {
	Body         []byte      `json:"body"`
	Headers      http.Header `json:"headers"`
	ResponseCode int         `json:"response_code"`
}

// Request sends a request using a connection
func Request(conn *Connection, requestType, url string, headers map[string]string) (responseBody []byte, responseCode int, responseHeaders http.Header, err error) {
	xlog.Debugf("Creating new %s request to: %s", requestType, url)

	// Create a basic request
	req, err := http.NewRequest(requestType, url, nil)
	if err != nil {
		xlog.Errorf("Error creating the HTTP request")
		return
	}

	// Add the headers
	for k, v := range headers {
		xlog.Debugf("Added header: {%s: %s}", k, v)
		req.Header.Add(k, v)
	}
	xlog.Debugf("Added %d headers", len(headers))

	// Perform the request
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		xlog.Errorf("Error performing HTTP %s request to %s", requestType, url)
		return
	}
	xlog.Debugf("Request made successfully with response: %+v", resp)

	// Check the return status code is valid
	responseCode = resp.StatusCode
	xlog.Debugf("Got response code: %d", responseCode)
	switch responseCode {
	case 200:
		break
	default:
		return nil, 0, nil, fmt.Errorf("Invalid reponse from Twitch (Error code: %d)", resp.StatusCode)
	}

	// Get the reponse headers
	responseHeaders = resp.Header

	// Close the body when we're done with it
	defer resp.Body.Close()

	// Get the body contents
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		xlog.Errorf("Could not read response body")
		return
	}
	xlog.Debugf("Got response body as: %s", responseBody)

	return
}
