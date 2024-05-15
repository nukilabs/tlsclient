package bandwidth

type Tracker interface {
	// AddWriteBytes adds the number of bytes written to the tracker.
	AddWriteBytes(n int64)
	// AddReadBytes adds the number of bytes read to the tracker.
	AddReadBytes(n int64)
}

type NoopTracker struct{}

func (NoopTracker) AddWriteBytes(n int64) {}
func (NoopTracker) AddReadBytes(n int64)  {}

func NewNoopTracker() Tracker {
	return &NoopTracker{}
}
