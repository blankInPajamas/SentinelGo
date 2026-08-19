package detector

import (
	"context"

	"github.com/blankInPajamas/SentinelGo/internal/alert"
)

type Detector interface {
	Start(ctx context.Context, alertChan chan<- alert.Alert) error
	Close() error
}
