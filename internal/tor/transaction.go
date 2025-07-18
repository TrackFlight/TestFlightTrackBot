package tor

import (
	"github.com/cretz/bine/tor"
	"math"
	"os"
)

type RequestTransaction struct {
	currentInstanceIndex int
	client               *Client
}

func (c *Client) NewTransaction(numRequests int) (*RequestTransaction, error) {
	c.mutex.Lock()
	neededInstances := int(math.Max(math.Floor(float64(numRequests/MaxRequestsPerInstance)), 1))
	deltaInstances := neededInstances - len(c.instances)
	if deltaInstances > 0 {
		for i := 0; i < deltaInstances; i++ {
			instance, err := newInstance()
			if err != nil {
				c.mutex.Unlock()
				return nil, err
			}
			c.instances = append(c.instances, instance)
		}
	} else if deltaInstances < 0 {
		removingInstances := c.instances[neededInstances:]
		for _, instance := range removingInstances {
			if err := instance.client.Close(); err != nil {
				c.mutex.Unlock()
				return nil, err
			}
			if err := os.RemoveAll(instance.client.DataDir); err != nil {
				c.mutex.Unlock()
				return nil, err
			}
		}
		c.instances = c.instances[:neededInstances]
	}
	return &RequestTransaction{
		client: c,
	}, nil
}

func (c *RequestTransaction) Close() {
	c.client.mutex.Unlock()
}

func (c *RequestTransaction) pickTorDialer() *tor.Dialer {
	c.currentInstanceIndex = (c.currentInstanceIndex + 1) % len(c.client.instances)
	return c.client.instances[c.currentInstanceIndex].dialer
}
