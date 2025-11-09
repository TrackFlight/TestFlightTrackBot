package tor

import (
	"fmt"
	"math/rand"
	stdHttp "net/http"
	"time"

	"github.com/Laky-64/http"
	"github.com/Laky-64/http/types"
)

func (c *RequestTransaction) ExecuteRequest(uri string) (*types.HTTPResult, error) {
	return http.ExecuteRequest(
		fmt.Sprintf("%s?nocache=%d", uri, time.Now().UnixNano()),
		http.Transport(
			&stdHttp.Transport{
				DialContext: c.pickTorDialer().DialContext,
			},
		),
		http.Headers(map[string]string{
			"User-Agent":      c.client.userAgents[rand.Intn(len(c.client.userAgents))],
			"Accept-Language": "en-US,en;q=0.9",
		}),
		http.Timeout(time.Second*5),
	)
}
