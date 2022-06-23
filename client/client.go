// SPDX-FileCopyrightText: 2021 Eric Neidhardt
// SPDX-License-Identifier: MIT

// Package client provides custom http client that measure throughput.
package client

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Request configures http request.
type Request struct {
	URL string

	PostBody    []byte
	ContentType string

	KeepAlive         bool
	AdditionalHeaders map[string]string
}

// Statistic contains measurement results.
type Statistic struct {
	// Overall number of bytes read.
	ReadThroughput int64
	// Overall number of bytes written.
	WriteThroughput int64

	// Overall number of performed requests, always equal to the sum of failed and successfull requests.
	RequestCount int
	SuccessCount int
	// Number of requests that failed with status code != 2xx.
	FailureCount int
	// Number of request that failed with error != nil while performing the request.
	NetworkFailedCount int
	// Number of request that failed with error != nil while reading response.
	IOFailedCount int
}

// Client is a custom http client that performs a request and collects measurements.
type Client struct {
	Statistic  Statistic
	Request    Request
	HTTPClient http.Client
}

// NewRequest creates a new request.
func NewRequest(
	url string,
	postDataFilePath string,
	postBody string,
	contentType string,
	keepAlive bool,
	authHeader string,
	additionalHeaders string,
) *Request {
	// preparing request
	var request = Request{
		URL:       url,
		KeepAlive: keepAlive,
	}

	// read optional post body
	if postDataFilePath != "" {
		data, err := os.ReadFile(postDataFilePath)
		if err != nil {
			log.Fatalf("Error while reading post body from file: %s Error: %s", postDataFilePath, err)
		}
		request.PostBody = data
	} else if postBody != "" {
		request.PostBody = []byte(postBody)
	}

	// create headers
	request.AdditionalHeaders = make(map[string]string)
	if contentType != "" {
		request.ContentType = contentType
	}
	if authHeader != "" {
		request.AdditionalHeaders["Authorization"] = authHeader
	}
	headers := strings.Split(additionalHeaders, ",")
	for _, header := range headers {
		keyValue := strings.Split(header, "=")
		if len(keyValue) == 2 {
			request.AdditionalHeaders[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		}
	}

	return &request
}

// NewClient creates a new client instance.
func NewClient(timeout time.Duration, Request Request) *Client {
	return &Client{
		HTTPClient: http.Client{
			Timeout: timeout,
		},
		Request: Request,
	}
}

// RunForDuration instructs the client to perform its request as often as possible for a given duration.
func (c *Client) RunForDuration(timeout time.Duration) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for time.Since(startTime) < timeout {
		c.PerformRequestWithContent(ctx)
	}
	// the last request can be interrupted by context timeout
	// we remove this from statistic
	if c.Statistic.NetworkFailedCount == 1 {
		c.Statistic.RequestCount--
		c.Statistic.NetworkFailedCount--
	}
}

// RunForAmount instructs the client to perform its request until a certain request count is reached.
func (c *Client) RunForAmount(requestCount int) {
	for i := 0; i < requestCount; i++ {
		c.PerformRequest()
	}
}

// PerformRequest instructs the client to perform its request once.
func (c *Client) PerformRequest() {
	c.PerformRequestWithContent(context.Background())
}

// PerformRequestWithContent instructs the client to perform its request once with a given context.
func (c *Client) PerformRequestWithContent(ctx context.Context) {
	// prepare request from configuration
	var req *http.Request
	var err error
	if c.Request.PostBody != nil {
		req, err = http.NewRequestWithContext(ctx, "POST", c.Request.URL, bytes.NewReader(c.Request.PostBody))
		req.Header.Set("Content-Type", c.Request.ContentType)
	} else {
		req, err = http.NewRequestWithContext(ctx, "GET", c.Request.URL, nil)
	}
	if err != nil {
		panic("Could not create http request")
	}

	if c.Request.KeepAlive {
		req.Header.Set("Connection", "keep-alive")
	} else {
		req.Header.Set("Connection", "close")
	}

	if c.Request.AdditionalHeaders != nil {
		for k, v := range c.Request.AdditionalHeaders {
			req.Header.Set(k, v)
		}
	}

	// perform request
	c.Statistic.RequestCount++
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.Statistic.NetworkFailedCount++
		return
	}
	defer resp.Body.Close()

	// write statistic
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode <= 299:
		c.Statistic.SuccessCount++
	default:
		c.Statistic.FailureCount++
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Statistic.IOFailedCount++
	}
	c.Statistic.ReadThroughput += int64(len(body))
	c.Statistic.WriteThroughput += int64(len(c.Request.PostBody))
}
