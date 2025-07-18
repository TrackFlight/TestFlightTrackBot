package tor

func (c *Client) Refresh() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, t := range c.instances {
		err := t.client.Control.Signal("NEWNYM")
		if err != nil {
			return err
		}
	}
	return nil
}
