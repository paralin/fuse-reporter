package remote

import (
	"encoding/json"
	"errors"
	"hash/crc32"
)

func (c *RequestContext) Validate() error {
	if c == nil {
		return errors.New("Request context cannot be nil.")
	}
	if c.HostIdentifier == "" {
		return errors.New("Host identifier cannot be empty.")
	}
	return nil
}

func (c *RemoteStreamConfig) ComputeCrc32() uint32 {
	str, _ := json.Marshal(&c.Streams)
	return crc32.ChecksumIEEE(str)
}

func (c *RemoteStreamConfig) FillCrc32() {
	c.Crc32 = int32(c.ComputeCrc32())
}
