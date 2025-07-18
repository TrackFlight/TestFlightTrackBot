package tor

import (
	"sync"
)

type Client struct {
	mutex     sync.Mutex
	instances []*Instance
}

func (c *Client) InstanceCount() int {
	return len(c.instances)
}
