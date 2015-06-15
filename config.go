package zipkin

import (
	"log"
	"os"
	"time"
)

// Config struct stores the configuration used to create a new ZipKin server. All fields
// are optional. Comment following each field indicates the default value that will be
// used.
type Config struct {
	ListenAddress       string        // "0.0.0.0:9410"
	InputBufferSize     int           // 32
	OutputBufferSize    int           // 32
	MaxConcurrentTraces int           // 128
	TraceTimeout        time.Duration // 3 * time.Second

	Logger *log.Logger // log.New(os.Stderr)
}

func fillDefaultConfig(c *Config) *Config {
	if c == nil {
		c = &Config{}
	}

	if c.ListenAddress == "" {
		c.ListenAddress = "0.0.0.0:9410"
	}
	if c.MaxConcurrentTraces == 0 {
		c.MaxConcurrentTraces = 128
	}
	if c.TraceTimeout == 0 {
		c.TraceTimeout = 3 * time.Second
	}

	if c.Logger == nil {
		// same as `var std` in stdlib src/log/log.go
		c.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}

	return c
}
