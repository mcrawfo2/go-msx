// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package httpclient

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/retry"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
	"strconv"
	"time"
)

const HeaderRetryAfter = "Retry-After"

// Retries with linear delay and Jitter (low random) (1, 2.452, 3.571, 4.357)
var DefaultHTTPClientRetryConfig = retry.RetryConfig{
	Attempts: 5,
	Delay:    1000,
	BackOff:  1.0,
	Linear:   true,
	Jitter:   1000,
}

type RetryPolicy = retryablehttp.CheckRetry

// NewRetryable returns an HTTP client that will retry requests based on the supplied
// retryPolicy and retry/backoff configuration.
// parses Retry-After on 429 response header by default
func NewRetryable(ctx context.Context, customizer Configurer, retryConfig retry.RetryConfig, backoff retryablehttp.Backoff, retryPolicy RetryPolicy) (*http.Client, error) {
	// Configure the go-msx http client
	client, err := New(ctx, customizer)
	if err != nil {
		return nil, err
	}

	// Set the default retry retryPolicy
	if retryPolicy == nil {
		retryPolicy = retryablehttp.ErrorPropagatedRetryPolicy
	}

	if backoff == nil {
		// Create a Retry instance to calculate backoff
		r := retry.NewRetry(ctx, retryConfig)

		backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			// tries to parse Retry-After response header when a http.StatusTooManyRequests
			// (HTTP Code 429) is found in the resp parameter
			// thank me (or hashicorp) later Meraki
			sleep := Parse429(resp)
			if sleep.IsPresent() {
				return sleep.Value() // should be in time.Duration already
			}

			return time.Duration(r.GetCurrentDelay(attemptNum + 1))
		}
	}

	// Create the retryable http client
	rhc := &retryablehttp.Client{
		HTTPClient: client,
		CheckRetry: retryPolicy,
		Backoff:    backoff,
		RetryMax:   retryConfig.Attempts - 1, // first attempt is not counted as retry
		Logger:     logger,                   // will be picked up as non leveled (just Printf)
	}

	// Return an actual http.Client
	return rhc.StandardClient(), nil
}

func Parse429(resp *http.Response) (result types.Optional[time.Duration]) {
	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
			if s, ok := resp.Header[HeaderRetryAfter]; ok {
				if sleep, err := strconv.ParseInt(s[0], 10, 64); err == nil {
					return types.OptionalOf(time.Duration(sleep) * time.Second)
				}
			}
		}
	}

	return types.OptionalEmpty[time.Duration]()
}
