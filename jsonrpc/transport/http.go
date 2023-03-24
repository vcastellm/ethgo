package transport

import (
	"encoding/json"

	"github.com/umbracle/ethgo/jsonrpc/codec"
	"github.com/valyala/fasthttp"
)

// HTTP is an http transport
type HTTP struct {
	addr    string
	client  *fasthttp.Client
	headers map[string]string
}

func newHTTP(addr string, headers map[string]string) *HTTP {
	return &HTTP{
		addr:    addr,
		client:  &fasthttp.Client{},
		headers: headers,
	}
}

// Close implements the transport interface
func (h *HTTP) Close() error {
	return nil
}

func (h *HTTP) Batch(requests []*codec.Request) ([]codec.Response, error) {
	raw, err := json.Marshal(requests)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	req.SetBody(raw)

	if err := h.client.Do(req, res); err != nil {
		return nil, err
	}

	// Decode json-rpc response
	var responses []codec.Response
	if err := json.Unmarshal(res.Body(), &responses); err != nil {
		return nil, err
	}

	return responses, nil
}

// Call implements the transport interface
func (h *HTTP) Call(method string, out interface{}, params ...interface{}) error {
	// Encode json-rpc request
	request, err := codec.NewRequest(method, params...)
	if err != nil {
		return err
	}
	raw, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(h.addr)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	for k, v := range h.headers {
		req.Header.Add(k, v)
	}
	req.SetBody(raw)

	if err := h.client.Do(req, res); err != nil {
		return err
	}

	// Decode json-rpc response
	var response codec.Response
	if err := json.Unmarshal(res.Body(), &response); err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error
	}

	if err := json.Unmarshal(response.Result, out); err != nil {
		return err
	}
	return nil
}

// SetMaxConnsPerHost sets the maximum number of connections that can be established with a host
func (h *HTTP) SetMaxConnsPerHost(count int) {
	h.client.MaxConnsPerHost = count
}
