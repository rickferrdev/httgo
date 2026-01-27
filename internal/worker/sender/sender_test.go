package sender

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/rickferrdev/httgo/internal/statistic/metrics"
)

func TestSender(t *testing.T) {
	var wgSender sync.WaitGroup
	const bf int = 5
	ch := make(chan metrics.Metrics, bf)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	request, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Error(err)
	}

	for range bf {
		wgSender.Go(func() {
			Sender(t.Context(), ch, http.DefaultClient, request)
		})
	}

	go func() {
		wgSender.Wait()
		close(ch)
	}()

	var receive int
	for range bf {
		select {
		case value := <-ch:
			receive++
			if value.Error != nil || value.StatusCode != 200 {
				t.Errorf("Request failed: err=%v, status=%d", value.Error, value.StatusCode)
			}
		case <-t.Context().Done():
			t.Error("closed context")
		}
	}

	if receive != bf {
		t.Errorf("Expected %d metrics, but got %d", bf, receive)
	}
}
