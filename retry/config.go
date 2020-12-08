package retry

type RetryConfig struct {
	Attempts int     `config:"default=3"`    // Maximum number of attempts
	Delay    int     `config:"default=500"`  // Milliseconds to wait between retries
	BackOff  float64 `config:"default=0.0"`  // Backoff factor
	Linear   bool    `config:"default=true"` // Use original delay when calculating current delay
}
