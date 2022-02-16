package flaresolverr

import (
	"io"
	"net/http"
	"strings"
)

func NewTransport(client *Client) *roundTripper {
	return &roundTripper{
		client: client,
	}
}

type roundTripper struct {
	client *Client
}

func (r roundTripper) RoundTrip(request *http.Request) (resp *http.Response, err error) {

	var rr *RequestResponse

	if request.Method == http.MethodGet {
		ops := RequestGetOps{}
		ops.URL = request.URL.String()
		rr, err = r.client.RequestGet(ops)
	} else {
		ops := RequestPostOps{}
		ops.URL = request.URL.String()
		rr, err = r.client.RequestPost(ops)
	}

	resp = &http.Response{}
	resp.StatusCode = rr.Solution.Status
	resp.Status = http.StatusText(rr.Solution.Status)
	resp.Body = io.NopCloser(strings.NewReader(rr.Solution.Response))

	for k, v := range rr.Solution.Headers {
		resp.Header.Set(k, v)
	}

	resp.ContentLength = -1
	resp.Uncompressed = true
	resp.Request = request

	return resp, nil
}
