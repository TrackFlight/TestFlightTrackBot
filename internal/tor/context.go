package tor

import (
	"sync"
)

type Client struct {
	mutex      sync.Mutex
	userAgents []string
	instances  []*Instance
}

func (c *Client) InstanceCount() int {
	return len(c.instances)
}
