package flaresolverr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

const (
	ResponseStatusOK    = "ok"
	ResponseStatusError = "error"
)

type Client struct {
	proto          string
	hostname       string
	port           int
	httpClient     *http.Client
	limiter        *rate.Limiter
	limiterContext context.Context
}

func NewClient(opts ...Option) *Client {

	c := &Client{}
	c.httpClient = http.DefaultClient
	c.hostname = "localhost"
	c.proto = "http"
	c.port = 8191

	envs := []Option{
		WithProtocol(os.Getenv("FSG_PROTO")),
		WithHostName(os.Getenv("FSG_HOST")),
		WithPortString(os.Getenv("FSG_PORT")),
	}

	opts = append(envs, opts...)

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type SessionCreateResponse struct {
	ResponseBase
	Session string `json:"session"`
}

func (c Client) SessionCreate(id string) (resp *SessionCreateResponse, err error) {

	v := map[string]string{
		"cmd":     "sessions.create",
		"session": id,
	}

	resp = &SessionCreateResponse{}
	err = c.request(v, resp)

	return resp, err
}

func (c Client) SessionDestroy(id string) (resp *ResponseBase, err error) {

	v := map[string]string{
		"cmd":     "sessions.destroy",
		"session": id,
	}

	resp = &ResponseBase{}
	err = c.request(v, resp)

	return resp, err
}

type SessionListResponse struct {
	Sessions []string `json:"sessions"`
}

func (c Client) SessionList() ([]string, error) {

	v := map[string]string{
		"cmd": "sessions.list",
	}

	sessions := &SessionListResponse{}
	err := c.request(v, sessions)

	return sessions.Sessions, err
}

type RequestGetOps struct {
	Command           string `json:"cmd"` // Ignore
	URL               string `json:"url,omitempty"`
	Session           string `json:"session,omitempty"`
	Timeout           int    `json:"maxTimeout,omitempty"` // milliseconds
	Cookies           string `json:"cookies,omitempty"`
	ReturnOnlyCookies string `json:"returnOnlyCookies,omitempty"`
	Proxy             struct {
		URL string `json:"url,omitempty"`
	} `json:"proxy,omitempty"`
}

type RequestResponse struct {
	ResponseBase
	Solution Solution `json:"solution"`
}

func (c Client) RequestGet(ops RequestGetOps) (*RequestResponse, error) {

	if c.limiter != nil {
		if err := c.limiter.Wait(c.limiterContext); err != nil {
			return nil, err
		}
	}

	ops.Command = "request.get"

	resp := &RequestResponse{}
	err := c.request(ops, resp)

	return resp, err
}

type RequestPostOps struct {
	RequestGetOps
	PostData string `json:"postData,omitempty"`
}

func (c Client) RequestPost(ops RequestPostOps) (*RequestResponse, error) {

	if c.limiter != nil {
		if err := c.limiter.Wait(c.limiterContext); err != nil {
			return nil, err
		}
	}

	ops.Command = "request.post"

	resp := &RequestResponse{}
	err := c.request(ops, resp)

	return resp, err
}

func (c Client) request(body interface{}, str interface{}) error {

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s://%s:%d/v1", c.proto, c.hostname, c.port), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, str)
}

type ResponseBase struct {
	Status         string `json:"status"`
	Message        string `json:"message"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Version        string `json:"version"`
}

type Solution struct {
	Url       string            `json:"url"`
	Status    int               `json:"status"`
	Headers   map[string]string `json:"headers"`
	Response  string            `json:"response"`
	Cookies   []Cookie          `json:"cookies"`
	UserAgent string            `json:"userAgent"`
}

type Cookie struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	Size     int     `json:"size"`
	HttpOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	Session  bool    `json:"session"`
	SameSite string  `json:"sameSite"`
}
