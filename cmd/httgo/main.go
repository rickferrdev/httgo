package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/rickferrdev/httgo/internal/commands/args"
	"github.com/rickferrdev/httgo/internal/statistic/metrics"
	"github.com/rickferrdev/httgo/internal/worker/consumer"
	"github.com/rickferrdev/httgo/internal/worker/sender"
)

var wgSender sync.WaitGroup
var wgConsumer sync.WaitGroup

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parse := args.Parse()
	ch := make(chan metrics.Metrics, parse.Goroutines)
	client := http.DefaultClient
	request, err := http.NewRequestWithContext(ctx, parse.MethodHttp, parse.Url, nil)
	if err != nil {
		log.Printf("Failed to instantiate a request: %v", err)
	}

	wgConsumer.Add(1)
	go func() {
		defer wgConsumer.Done()
		consumer.Orchestrator(ctx, ch)
	}()

	for range parse.Goroutines {
		wgSender.Go(func() {
			sender.Sender(ctx, ch, client, request)
		})
	}

	wgSender.Wait()
	close(ch)
	wgConsumer.Wait()
}
