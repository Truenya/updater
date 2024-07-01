package util

import (
	"fmt"
	"time"
)

type Progress struct {
	Message      string
	LastStageDur time.Duration
	Elapsed      time.Duration
	Est          time.Duration
	Error        error
}

func (p *Progress) String() string {
	if p.Error != nil {
		return p.Error.Error()
	}

	return fmt.Sprintf(
		"%s | last stage: %s \tfull: %s \test: %s",
		p.Message,
		p.LastStageDur.Round(time.Millisecond),
		p.Elapsed.Round(time.Millisecond),
		p.Est.Round(time.Millisecond),
	)
}
