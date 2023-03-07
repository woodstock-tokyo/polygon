package polygon

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const apiURL = "https://api.polygon.io/v2" // use v2 as default

// Client models a client to consume the Polygon Cloud API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	launchPad  bool
	Edge       *Edge
}

// Edge is used for polygon launchpad
type Edge struct {
	id        string
	ipAddress string
}

// polygon api versions are not unified, sometimes we have to switch to v1
func (c Client) UseV1Endpoints() Client {
	c.baseURL = strings.Replace(c.baseURL, "v2", "v1", 1)
	return c
}

// polygon api versions are not unified, sometimes we have to switch to v3
func (c Client) UseV3Endpoints() Client {
	c.baseURL = strings.Replace(c.baseURL, "v2", "v3", 1)
	return c
}

// polygon api versions are not unified, sometimes we have to switch to vX
func (c Client) UseVXEndpoints() Client {
	c.baseURL = strings.Replace(c.baseURL, "v2", "vX", 1)
	return c
}

// Error represents an Polygon API error
type Error struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error"`
	RequestID    string `json:"request_id"`
}

// ClientOption applies an option to the client.
type ClientOption func(*Client)

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s (request id: %s)", e.Status, e.ErrorMessage, e.RequestID)
}

// NewClient creates a client with the given authorization token.
func NewClient(token string, options ...ClientOption) *Client {
	client := &Client{
		token:      token,
		httpClient: &http.Client{Timeout: time.Second * 60},
	}

	// apply options
	for _, applyOption := range options {
		applyOption(client)
	}

	// set default values
	if client.baseURL == "" {
		client.baseURL = apiURL
	}

	return client
}

// WithHTTPClient sets the http.Client for a new Polygon Client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// WithSecureHTTPClient sets a secure http.Client for a new Polygon Client
func WithSecureHTTPClient() ClientOption {
	return func(client *Client) {
		client.httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}}
	}
}

// WithBaseURL sets the baseURL for a new Polygon Client
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		client.baseURL = baseURL
	}
}

// WithEdge set edge for Polygon launchpad
func WithEdge(id, ipAddress string) ClientOption {
	return func(client *Client) {
		client.Edge = &Edge{
			id:        id,
			ipAddress: ipAddress,
		}

		client.launchPad = true
	}
}

// GetJSON gets the JSON data from the given endpoint.
func (c *Client) GetJSON(ctx context.Context, endpoint string, v any) error {
	u, err := c.url(endpoint, map[string]string{"apiKey": c.token})
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// GetJSONWithQueryParams gets the JSON data from the given endpoint with the query parameters attached.
func (c *Client) GetJSONWithQueryParams(ctx context.Context, endpoint string, queryParams map[string]string, v any) error {
	queryParams["apiKey"] = c.token
	u, err := c.url(endpoint, queryParams)
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// Fetches JSON content from the given URL and unmarshals it into `v`.
func (c *Client) FetchURLToJSON(ctx context.Context, u *url.URL, v any) error {
	data, err := c.getBytes(ctx, u.String())
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// GetJSONWithoutToken gets the JSON data from the given endpoint without
// adding a token to the URL.
func (c *Client) GetJSONWithoutToken(ctx context.Context, endpoint string, v any) error {
	u, err := c.url(endpoint, nil)
	if err != nil {
		return err
	}
	return c.FetchURLToJSON(ctx, u, v)
}

// GetBytes gets the data from the given endpoint.
func (c *Client) GetBytes(ctx context.Context, endpoint string) ([]byte, error) {
	u, err := c.url(endpoint, map[string]string{"apiKey": c.token})
	if err != nil {
		return nil, err
	}
	return c.getBytes(ctx, u.String())
}

// GetFloat64 gets the number from the given endpoint.
func (c *Client) GetFloat64(ctx context.Context, endpoint string) (float64, error) {
	b, err := c.GetBytes(ctx, endpoint)
	if err != nil {
		return 0.0, err
	}
	return strconv.ParseFloat(string(b), 64)
}

func (c *Client) getBytes(ctx context.Context, address string) ([]byte, error) {
	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return []byte{}, err
	}

	if c.Edge != nil {
		req.Header.Set("X-Polygon-Edge-ID", c.Edge.id)
		req.Header.Set("X-Polygon-Edge-IP-Address", c.Edge.ipAddress)
	}

	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	// Even if GET didn't return an error, check the status code to make sure
	// everything was ok.
	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		msg := ""

		if err == nil {
			msg = string(b)
		}

		return []byte{}, Error{Status: resp.Status, ErrorMessage: msg}
	}
	return io.ReadAll(resp.Body)
}

// Returns an URL object that points to the endpoint with optional query parameters.
func (c *Client) url(endpoint string, queryParams map[string]string) (*url.URL, error) {
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if queryParams != nil {
		q := u.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u, nil
}

func (c Client) endpointWithOpts(endpoint string, opts any) (string, error) {
	if opts == nil {
		return endpoint, nil
	}
	v, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	optParams := v.Encode()
	if optParams != "" {
		endpoint = fmt.Sprintf("%s?%s", endpoint, optParams)
	}

	return endpoint, nil
}

// ttoa time to string helper
func ttoa(t time.Time, layout ...string) string {
	l := "2006-01-02"
	if len(layout) > 0 {
		l = layout[0]
	}
	return t.Format(l)
}
