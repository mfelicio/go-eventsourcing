package sourcing

import ()

type EventStore interface {
	Load(aggregateId string) ([]*Event, error)
	LoadSince(aggregateId string, start int) ([]*Event, error)
	Append(aggregateId string, expectedVersion int, events []*Event) error
}

/**

Simple Memory EventStore implementation

*/

func NewMemoryEventStore() EventStore {
	return &memoryEventStore{data: make(map[string][]*Event)}
}

type memoryEventStore struct {
	nextId string
	data   map[string][]*Event
}

func (this *memoryEventStore) Load(aggregateId string) ([]*Event, error) {

	if events, ok := this.data[aggregateId]; ok {
		return events, nil
	}

	return []*Event{}, nil
}

func (this *memoryEventStore) LoadSince(aggregateId string, start int) ([]*Event, error) {

	if events, ok := this.data[aggregateId]; ok {
		return events[start:], nil
	}

	return []*Event{}, nil
}

func (this *memoryEventStore) Append(aggregateId string, expectedVersion int, events []*Event) error {

	if _, ok := this.data[aggregateId]; ok {
		for _, ev := range events {
			this.data[aggregateId] = append(this.data[aggregateId], ev)
		}
	} else {
		this.data[aggregateId] = events
	}

	return nil

}
