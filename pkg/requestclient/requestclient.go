package requestclient

import (
	"context"
	"encoding/json"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/valyala/fasthttp"
)

type Request struct {
	URL     string
	Method  string
	Body    interface{}
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
}

type IRequestClient interface {
	HandleRequest(ctx context.Context, req Request) (*Response, error)
	GetReq(ctx context.Context, headers map[string]string, url string) (*Response, error)
}

type requestClient struct{}

func NewClientRequest() IRequestClient {
	return &requestClient{}
}

func (r *requestClient) HandleRequest(ctx context.Context, req Request) (*Response, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	var externalSegment *newrelic.ExternalSegment
	if tx := newrelic.FromContext(ctx); tx != nil {
		externalSegment = &newrelic.ExternalSegment{URL: req.URL}
		externalSegment.StartTime = tx.StartSegmentNow()
		defer externalSegment.End()
	}

	request.SetRequestURI(req.URL)
	if req.Body != nil {
		body, encodeErr := json.Marshal(req.Body)
		if encodeErr != nil {
			return nil, encodeErr
		}
		request.SetBody(body)
	}

	request.Header.SetMethod(req.Method)
	for key, header := range req.Headers {
		request.Header.Set(key, header)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := fasthttp.Do(request, resp)
	if err != nil {
		return nil, err
	}

	if externalSegment != nil {
		externalSegment.SetStatusCode(resp.StatusCode())
	}

	return &Response{
		Body:       resp.Body(),
		StatusCode: resp.StatusCode(),
	}, nil
}

func (r *requestClient) GetReq(ctx context.Context, headers map[string]string, url string) (*Response, error) {
	return r.HandleRequest(ctx, Request{
		URL:     url,
		Method:  fasthttp.MethodGet,
		Headers: headers,
	})
}
