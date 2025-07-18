package testflight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/Laky-64/http"
	"math/rand"
)

func loadUserAgents() ([]string, error) {
	request, err := http.ExecuteRequest(UserAgentListURL)
	if err != nil {
		return nil, err
	}

	var userAgents []string
	scanner := bufio.NewScanner(bytes.NewReader(request.Body))
	for scanner.Scan() {
		var userAgent UserAgent
		if err = json.Unmarshal(scanner.Bytes(), &userAgent); err == nil {
			userAgents = append(userAgents, userAgent.UserAgent)
		}
	}
	return userAgents, scanner.Err()
}

func pickRandomUserAgent(userAgents []string) string {
	return userAgents[rand.Intn(len(userAgents))]
}
