package tor

import (
	"github.com/Laky-64/http"
	"github.com/Laky-64/http/types"
	stdHttp "net/http"
)

func (c *RequestTransaction) ExecuteRequest(uri string, userAgent string) (*types.HTTPResult, error) {
	return http.ExecuteRequest(
		uri,
		http.Transport(
			&stdHttp.Transport{
				DialContext: c.pickTorDialer().DialContext,
			},
		),
		http.Headers(map[string]string{
			"User-Agent": userAgent,
		}),
	)
}
