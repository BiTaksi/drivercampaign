package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/valyala/fasthttp"

	"github.com/BiTaksi/drivercampaign/pkg/nrclient"
	"github.com/BiTaksi/drivercampaign/pkg/utils"
)

func getTransactionName(c *fiber.Ctx) string {
	return fmt.Sprintf("%s %s", c.Method(), c.Path())
}

func getTransportType(ctx *fasthttp.RequestCtx) newrelic.TransportType {
	if ctx.IsTLS() {
		return newrelic.TransportHTTPS
	}
	return newrelic.TransportHTTP
}

func getURL(ctx *fasthttp.RequestCtx) *url.URL {
	u := string(ctx.RequestURI())
	parse, _ := url.ParseRequestURI(u)
	return parse
}

func getHTTPHeader(ctx *fasthttp.RequestCtx) http.Header {
	h := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		h.Set(string(k), string(v))
	})
	return h
}

type ResponseWriter struct {
	Code      int
	HeaderMap http.Header
	Body      *bytes.Buffer
	ResCtx    *fasthttp.Response
}

func NewResponseWriter(resCtx *fasthttp.Response) *ResponseWriter {
	return &ResponseWriter{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      http.StatusOK,
		ResCtx:    resCtx,
	}
}

func (rw *ResponseWriter) Header() http.Header {
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
	return rw.HeaderMap
}

func (rw *ResponseWriter) writeHeader() {
	rw.ResCtx.Header.VisitAll(func(k, v []byte) {
		rw.HeaderMap.Set(string(k), string(v))
	})
}

func (rw *ResponseWriter) Write(buf []byte) (int, error) {
	rw.writeHeader()
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.Code = code
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
	rw.writeHeader()
}

func NewRelicMiddleware(nr nrclient.INewRelicInstance) func(c *fiber.Ctx) (err error) {
	return func(c *fiber.Ctx) (err error) {
		txn := nr.Application().StartTransaction(getTransactionName(c))
		txn.SetWebRequest(newrelic.WebRequest{
			Host:      string(c.Context().Host()),
			Header:    getHTTPHeader(c.Context()),
			Method:    c.Method(),
			URL:       getURL(c.Context()),
			Transport: getTransportType(c.Context()),
		})
		txn.AddAttribute(utils.RequestIDKey, c.Locals(utils.RequestIDKey))
		defer txn.End()

		rw := txn.SetWebResponse(NewResponseWriter(c.Response()))
		c.Locals(utils.NewrelicContextKey, newrelic.NewContext(c.Context(), txn))

		err = c.Next()
		if err == nil {
			rw.WriteHeader(c.Response().StatusCode())
			return err
		}

		errCode := http.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			errCode = e.Code
		}

		txn.SetWebResponse(nil).WriteHeader(errCode)
		return err
	}
}
