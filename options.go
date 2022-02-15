package flaresolverr

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Option func(*Client)

func WithProtocol(proto string) Option {
	return func(c *Client) {
		c.proto = proto
	}
}

func WithPort(port int) Option {
	return func(c *Client) {
		c.port = port
	}
}

func WithHostName(hostname string) Option {
	return func(c *Client) {
		c.hostname = hostname
	}
}

func WithClient(httpclient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpclient
	}
}

func WithRate(duration time.Duration, ctx context.Context) Option {
	return func(c *Client) {
		c.limiter = rate.NewLimiter(rate.Every(duration), 1)
		c.limiterContext = ctx
	}
}
