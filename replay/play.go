package replay

import (
	"time"
)

// Represents one request try
type Play struct {
	Bytes  int
	Writes int
	Status int
	Time   time.Duration

	start time.Time
	end   time.Time
}

func (p *Play) Start() {
	p.start = time.Now()
}

func (p *Play) Stop() {
	p.end = time.Now()
	p.Time = p.end.Sub(p.start)
}

func (p *Play) Finished() bool {
	return p.end.After(p.start) && p.Status != 0
}
