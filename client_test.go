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
	data, err := client.getBytes(ctx, "/v1/test-endpoint")
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
	data, err := client.getBytes(ctx, "/v1/test-endpoint")
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
