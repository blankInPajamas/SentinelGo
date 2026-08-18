package collector

import "context"

// Collector defines the contract for all log collectors.
type Collector interface {
	Start(ctx context.Context, handler func(line string)) error
	Close() error
}
