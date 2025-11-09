package tor

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/Laky-64/http"
)

func (c *Client) loadUserAgents() error {
	request, err := http.ExecuteRequest(UserAgentListURL)
	if err != nil {
		return err
	}

	var userAgents []string
	scanner := bufio.NewScanner(bytes.NewReader(request.Body))
	for scanner.Scan() {
		var userAgent UserAgent
		if err = json.Unmarshal(scanner.Bytes(), &userAgent); err == nil {
			userAgents = append(userAgents, userAgent.UserAgent)
		}
	}
	c.userAgents = userAgents
	return scanner.Err()
}
