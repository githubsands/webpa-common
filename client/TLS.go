package client

import "crypto/tls"

type tLSConfig struct {
	// InsecureSkipVerify if set to true will disable TLS host verification.
	InsecureSkipVerify bool
}

func newTLSConfig(c ClientConfig) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: c.InsecureSkipVerify(),
	}
}

func (o *tLSConfig) InsecureSkipVerify() {
	if o != nil && o.tLSConfig.InsecureSkipVerify != nil {
		return o.tLSConfig.InsecureSkipVerify
	}

	return nil
}
