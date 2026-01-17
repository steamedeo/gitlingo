package relay

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"steamedeo.dev/gitlingo/internal/config"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// RelayOptions configures an HTTP request.
type RelayOptions struct {
	Context        context.Context // request cancellation and timeouts
	URL            string          // target endpoint
	Method         string          // HTTP verb (GET, POST, etc.)
	RequestPayload io.Reader       // request payload (optional)
}

type Relay struct {
	Options RelayOptions
}

func NewRelay(o RelayOptions) (*Relay, error) {
	if o.Context == nil {
		return nil, fmt.Errorf("context parameter missing")
	}
	if o.URL == "" {
		return nil, fmt.Errorf("url parameter missing")
	}
	if o.Method == "" {
		return nil, fmt.Errorf("method parameter missing")
	}
	return &Relay{
		Options: o,
	}, nil
}

func (r *Relay) SendHTTPRequest() (*http.Response, error) {
	req, err := http.NewRequestWithContext(r.Options.Context, r.Options.Method, r.Options.URL, r.Options.RequestPayload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+config.AppConfig.GithubToken)
	req.Header.Add("User-Agent", "gitlingo")
	res, err := performHTTPWithBackoff(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func performHTTPWithBackoff(r *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	i := 0
	backoff := 1
	for i < 3 {
		res, err = httpClient.Do(r)
		if err == nil {
			if res.StatusCode >= 400 && res.StatusCode < 600 {
				return nil, fmt.Errorf("HTTP request failed with status %d", res.StatusCode)
			}
			return res, nil
		}
		time.Sleep(time.Duration(backoff) * time.Second)
		i++
		backoff += backoff
	}
	return res, err
}
