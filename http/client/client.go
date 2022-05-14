package client

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	URL "net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/trivelaapp/go-kit/errors"
)

type httpClientProvider interface {
	Do(request *http.Request) (*http.Response, error)
}

// Client provides methods for making REST requests.
type Client struct {
	http httpClientProvider
}

// New creates a new Client instance.
func New(timeout time.Duration) Client {
	return Client{
		http: &http.Client{
			Timeout:   timeout,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

// Patch execute a http PATCH method with application/json headers.
func (c Client) Patch(ctx context.Context, request HTTPRequest) (rst HTTPResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest(ctx, "PATCH", request)
}

// Put execute a http PUT method with application/json headers.
func (c Client) Put(ctx context.Context, request HTTPRequest) (rst HTTPResult, err error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest(ctx, "PUT", request)
}

// Post execute a http POST method with application/json headers.
func (c Client) Post(ctx context.Context, request HTTPRequest) (HTTPResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	request.Headers["content-type"] = "application/json"
	return c.processRequest(ctx, "POST", request)
}

// Delete execute a http DELETE method with application/json headers.
func (c Client) Delete(ctx context.Context, request HTTPRequest) (HTTPResult, error) {
	if request.Headers == nil {
		request.Headers = make(map[string]string)
	}
	return c.processRequest(ctx, "DELETE", request)
}

// Get execute a http GET method.
func (c Client) Get(ctx context.Context, request HTTPRequest) (HTTPResult, error) {
	return c.processRequest(ctx, "GET", request)
}

func (c Client) processRequest(ctx context.Context, method string, request HTTPRequest) (HTTPResult, error) {
	queryValues := URL.Values{}

	for key, value := range request.QueryParams {
		queryValues.Add(key, value)
	}

	url, err := URL.Parse(request.URL)
	if err != nil {
		return HTTPResult{}, errors.New("error on parsing the request url")
	}
	url.RawQuery = queryValues.Encode()

	httpRequest, err := http.NewRequestWithContext(ctx, method, url.String(), bytes.NewBuffer(request.Body))
	if err != nil {
		return HTTPResult{}, err
	}
	for key, value := range request.Headers {
		httpRequest.Header.Add(key, value)
	}

	res, err := c.http.Do(httpRequest)
	if err != nil {
		return HTTPResult{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return HTTPResult{}, err
	}

	location, _ := res.Location()

	return HTTPResult{
		Response:   body,
		Location:   location.String(),
		StatusCode: res.StatusCode,
	}, nil
}
