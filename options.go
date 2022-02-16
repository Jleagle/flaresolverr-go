package flaresolverr

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type Option func(*Client)

func WithProtocol(proto string) Option {
	return func(c *Client) {
		if proto != "" {
			c.proto = proto
		}
	}
}

func WithPort(port int) Option {
	return func(c *Client) {
		if port > 0 {
			c.port = port
		}
	}
}

func WithPortString(port string) Option {
	p, _ := strconv.Atoi(port)
	return WithPort(p)
}

func WithHostName(hostname string) Option {
	return func(c *Client) {
		if hostname != "" {
			c.hostname = hostname
		}
	}
}

func WithClient(httpclient *http.Client) Option {
	return func(c *Client) {
		if httpclient != nil {
			c.httpClient = httpclient
		}
	}
}

func WithRate(duration time.Duration, ctx context.Context) Option {
	return func(c *Client) {
		if duration > 0 {
			c.limiter = rate.NewLimiter(rate.Every(duration), 1)
			c.limiterContext = ctx
		}
	}
}
