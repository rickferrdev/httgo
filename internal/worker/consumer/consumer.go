package consumer

import (
	"context"
	"log"
	"time"

	"github.com/rickferrdev/httgo/internal/statistic/metrics"
)

func Orchestrator(ctx context.Context, ch <-chan metrics.Metrics) {
	var (
		total         int
		successful    int
		duration      time.Duration
		dnsTotal      time.Duration
		tcpTotal      time.Duration
		serverTotal   time.Duration
		transferTotal time.Duration
	)

	startWallClock := time.Now()

	for {
		select {
		case <-ctx.Done():
			log.Println("channel cancelled")
			return
		case value, ok := <-ch:
			if !ok {
				wallClock := time.Since(startWallClock)
				Render(total, successful, duration, wallClock)
				return
			}

			total++
			if value.Error == nil && value.StatusCode >= 200 && value.StatusCode <= 300 {
				successful++
			}

			duration += value.Duration
			dnsTotal += value.DNSLookup
			tcpTotal += value.TCPConnection
			serverTotal += value.ServerThinking
			transferTotal += value.TransferTime
		}
	}
}

func Render(total int, success int, totalDuration time.Duration, wallClock time.Duration) {
	avg := time.Duration(0)
	if total > 0 {
		avg = totalDuration / time.Duration(total)
	}
	log.Printf("Summary: Processed %d requests with a %d success count.", total, success)
	log.Printf("Performance: Total processing time was %v (Avg: %v per request).", totalDuration, avg)
	log.Printf("User Experience: Tasks completed in %v (wall clock).", wallClock)
}
