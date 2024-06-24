package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/valyala/fasthttp"
)

const (
	clientName     = "DriverPayment/1.0"
	defaultTimeout = 10 * time.Second
	maxConnPerHost = 1000
	dialTimeout    = 7 * time.Second
)

type Request struct {
	URL     string
	Method  string
	Body    interface{}
	Headers map[string]string

	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Body       []byte
}

type IHTTPClient interface {
	HandleRequest(ctx context.Context, req Request) (*Response, error)
	IsSuccessStatusCode(resp *Response) bool
	GetJSONHeaders() map[string]string
	HandleException(resp *Response) error
	HandleInternalException(resp *Response) error
}

type httpClient struct {
	client *fasthttp.Client
}

func NewHTTPClient() IHTTPClient {
	fc := &fasthttp.Client{
		Name:                     clientName,
		MaxConnsPerHost:          maxConnPerHost,
		NoDefaultUserAgentHeader: true,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, dialTimeout)
		},
	}

	return &httpClient{client: fc}
}

func (h *httpClient) HandleRequest(ctx context.Context, req Request) (*Response, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	var externalSegment *newrelic.ExternalSegment
	if tx := newrelic.FromContext(ctx); tx != nil {
		externalSegment = &newrelic.ExternalSegment{URL: req.URL}
		externalSegment.StartTime = tx.StartSegmentNow()
		defer externalSegment.End()
	}

	request.SetRequestURI(req.URL)
	if req.Body != nil {
		body, encodeErr := h.prepareBody(req.Body)
		if encodeErr != nil {
			return nil, encodeErr
		}
		request.SetBody(body)
	}

	request.Header.SetMethod(req.Method)
	for key, header := range req.Headers {
		request.Header.Set(key, header)
	}

	if req.Timeout <= 0 {
		req.Timeout = defaultTimeout
	}

	if err := h.client.DoTimeout(request, resp, req.Timeout); err != nil {
		return nil, fmt.Errorf("request err: %v", err)
	}

	respBody := resp.Body()
	respStatusCode := resp.StatusCode()

	var bytes []byte
	bytes = append(bytes, respBody...)

	if externalSegment != nil {
		externalSegment.SetStatusCode(respStatusCode)
	}

	return &Response{
		Body:       bytes,
		StatusCode: respStatusCode,
	}, nil
}

func (h *httpClient) IsSuccessStatusCode(resp *Response) bool {
	return resp.StatusCode >= fasthttp.StatusOK && resp.StatusCode < fasthttp.StatusMultipleChoices
}

func (h *httpClient) GetJSONHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"

	return headers
}

func (h *httpClient) prepareBody(b interface{}) ([]byte, error) {
	if byteBody, byteBodyOk := b.([]byte); byteBodyOk {
		return byteBody, nil
	}

	body, encodeErr := json.Marshal(b)
	if encodeErr != nil {
		return nil, encodeErr
	}

	return body, nil
}

func (h *httpClient) HandleException(resp *Response) error {
	respErr := ResponseErrorBag{
		Response: Response{
			StatusCode: resp.StatusCode,
			Body:       resp.Body,
		},
	}

	if jsonErr := json.Unmarshal(resp.Body, &respErr); jsonErr != nil {
		respErr.Cause = fmt.Errorf("json err: %v", jsonErr)
	} else {
		respErr.Cause = errors.New("cause: nil")
	}

	return respErr
}

func (h *httpClient) HandleInternalException(resp *Response) error {
	respErr := ResponseInternalErrorBag{
		Response: Response{
			StatusCode: resp.StatusCode,
			Body:       resp.Body,
		},
	}

	if jsonErr := json.Unmarshal(resp.Body, &respErr); jsonErr != nil {
		respErr.Cause = fmt.Errorf("json err: %v", jsonErr)
	} else {
		respErr.Cause = errors.New("cause: nil")
	}

	return respErr
}
