// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package retry

type RetryConfig struct {
	Attempts int     `config:"default=3"`    // Maximum number of attempts
	Delay    int     `config:"default=500"`  // Milliseconds to wait between retries
	BackOff  float64 `config:"default=0.0"`  // Backoff factor
	Linear   bool    `config:"default=true"` // Use original delay when calculating current delay
	Jitter   int     `config:"default=0"`    // in milliseconds. linearBackoff would have random jitter from Delay to Delay + Jitter
}
