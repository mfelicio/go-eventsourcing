package sourcing

import (
	"time"
)

type Aggregate interface {
	Id() string

	//The original version when the source aggregate was loaded
	Version() int

	//New events since original version
	Changes() []*Event

	On(ev interface{})
}

func loadAggregate(binder Binder, id string, source interface{}, events []*Event) *aggregate {

	aggregate := &aggregate{
		id:      id,
		version: len(events),
		binder:  binder,
		source:  source,
		events:  events,
	}

	for _, ev := range events {
		//Apply event to source state
		binder.Apply(source, ev.Payload)
	}

	return aggregate
}

type aggregate struct {
	id      string
	version int
	events  []*Event
	binder  Binder
	source  interface{}
}

func (this *aggregate) Id() string {
	return this.id
}

func (this *aggregate) Version() int {
	return this.version
}

func (this *aggregate) Changes() []*Event {

	//TODO: only works because aggregates are being loaded from all events
	//When snapshot support is added the slice below will not work

	//returns all events since the source.version
	return this.events[this.version:]
}

func (this *aggregate) On(ev interface{}) {

	//Apply event to source state
	this.binder.Apply(this.source, ev)

	//create new event
	event := &Event{
		AggregateId: this.id,
		SequenceId:  len(this.events) - this.version + 1,
		Time:        time.Now(),
		Name:        this.binder.GetEventName(ev),
		Payload:     ev,
	}

	//Add event to aggregate events
	this.events = append(this.events, event)
}
