package remote

import "errors"

func (c *RequestContext) Validate() error {
	if c == nil {
		return errors.New("Request context cannot be nil.")
	}
	if c.HostIdentifier == "" {
		return errors.New("Host identifier cannot be empty.")
	}
	return nil
}
