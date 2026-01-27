package sender

import (
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/rickferrdev/httgo/internal/statistic/metrics"
)

func Sender(ctx context.Context, ch chan<- metrics.Metrics, client *http.Client, request *http.Request) {
	var (
		start                        = time.Now()
		dnsStart                     = time.Now()
		connStart                    = time.Now()
		dnsDone, connDone, firstByte time.Duration
	)

	trace := &httptrace.ClientTrace{
		DNSStart:             func(di httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:              func(di httptrace.DNSDoneInfo) { dnsDone = time.Since(dnsStart) },
		ConnectStart:         func(network, addr string) { connStart = time.Now() },
		ConnectDone:          func(network, addr string, err error) { connDone = time.Since(connStart) },
		GotFirstResponseByte: func() { firstByte = time.Since(start) },
	}

	request = request.WithContext(httptrace.WithClientTrace(ctx, trace))

	start = time.Now()
	response, err := client.Do(request)
	if err != nil {
		ch <- metrics.Metrics{
			Duration: time.Since(start),
			Error:    err,
		}
		return
	}
	defer response.Body.Close()
	io.Copy(io.Discard, response.Body)

	totalDuration := time.Since(start)
	ch <- metrics.Metrics{
		Duration:   totalDuration,
		StatusCode: response.StatusCode,

		DNSLookup:      dnsDone,
		TCPConnection:  connDone,
		ServerThinking: firstByte,
		TransferTime:   totalDuration - firstByte,
		Error:          nil,
	}
}
