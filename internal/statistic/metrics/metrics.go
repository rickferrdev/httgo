package metrics

import "time"

type Metrics struct {
	StatusCode int
	Duration   time.Duration

	DNSLookup      time.Duration
	TCPConnection  time.Duration
	TLSHandshake   time.Duration
	ServerThinking time.Duration
	TransferTime   time.Duration

	Error error
}
