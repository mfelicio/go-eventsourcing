package sourcing

import (
	"time"
)

type Event struct {
	AggregateId string
	SequenceId  int
	Time        time.Time
	Name        string
	Payload     interface{}
}
