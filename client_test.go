// Copyright (c) 2026 Woodstock K.K.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package polygon

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	calls int
	errs  []error
	resps []*http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	defer func() { m.calls++ }()
	if m.calls < len(m.errs) {
		return m.resps[m.calls], m.errs[m.calls]
	}
	return nil, errors.New("no more mock responses")
}

func TestGetBytesRetry(t *testing.T) {
	goawayErr := errors.New("http2: server sent GOAWAY and closed the connection; LastStreamID=1")
	body := io.NopCloser(bytes.NewBufferString(`{"status":"OK"}`))
	successResp := &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	rt := &mockRoundTripper{
		errs:  []error{goawayErr, nil},
		resps: []*http.Response{nil, successResp},
	}

	client := NewClient("test-api-key")
	client.httpClient.Transport = rt

	ctx := context.Background()
	data, err := client.getBytes(ctx, "/v1/test-endpoint", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(data) != `{"status":"OK"}` {
		t.Fatalf("expected response body to be %s, got %s", `{"status":"OK"}`, string(data))
	}
	if rt.calls != 2 {
		t.Fatalf("expected 2 calls to RoundTrip, got %d", rt.calls)
	}
}

func TestGetBytesNoRetrySuccess(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{"status":"OK"}`))
	successResp := &http.Response{
		StatusCode: 200,
		Body:       body,
	}

	rt := &mockRoundTripper{
		errs:  []error{nil},
		resps: []*http.Response{successResp},
	}

	client := NewClient("test-api-key")
	client.httpClient.Transport = rt

	ctx := context.Background()
	data, err := client.getBytes(ctx, "/v1/test-endpoint", true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(data) != `{"status":"OK"}` {
		t.Fatalf("expected response body to be %s, got %s", `{"status":"OK"}`, string(data))
	}
	if rt.calls != 1 {
		t.Fatalf("expected 1 call to RoundTrip, got %d", rt.calls)
	}
}
