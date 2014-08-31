package sourcing

import (
	"reflect"
)

type ApplyEventHandler func(source interface{}, ev interface{})
type CreateHandler func() interface{}
type InitHandler func(source interface{}, aggregate Aggregate)

type Binder interface {
	Apply(source interface{}, ev interface{})

	Create() interface{}
	Init(source interface{}, aggregate Aggregate)

	GetSourceName() string
	GetEventName(ev interface{}) string
}

type Binding interface {
	Binder

	On(eventType reflect.Type, handler ApplyEventHandler) Binding
	OnNamed(eventType reflect.Type, name string, handler ApplyEventHandler) Binding
}

type binding struct {
	sourceType    reflect.Type
	sourceName    string
	create        CreateHandler
	init          InitHandler
	eventBindings map[reflect.Type]*eventBinding
}

type eventBinding struct {
	applyEvent ApplyEventHandler
	eventName  string
}

func newBinding(sourceType reflect.Type, sourceName string, create CreateHandler, init InitHandler) *binding {

	return &binding{
		sourceType:    sourceType,
		sourceName:    sourceName,
		create:        create,
		init:          init,
		eventBindings: make(map[reflect.Type]*eventBinding),
	}
}

func newEventBinding(eventType reflect.Type, eventName string, handler ApplyEventHandler) *eventBinding {

	return &eventBinding{
		eventName:  eventName,
		applyEvent: handler,
	}
}

func (this *binding) On(eventType reflect.Type, handler ApplyEventHandler) Binding {

	return this.OnNamed(eventType, eventType.Name(), handler)
}

func (this *binding) OnNamed(eventType reflect.Type, name string, handler ApplyEventHandler) Binding {

	this.eventBindings[eventType] = newEventBinding(eventType, name, handler)
	return this
}

func (this *binding) GetSourceName() string {

	return this.sourceName
}

func (this *binding) GetEventName(ev interface{}) string {

	return this.eventBindings[reflect.TypeOf(ev)].eventName
}

func (this *binding) Create() interface{} {

	return this.create()
}

func (this *binding) Init(source interface{}, aggregate Aggregate) {

	this.init(source, aggregate)
}

func (this *binding) Apply(source interface{}, ev interface{}) {

	evBinding := this.eventBindings[reflect.TypeOf(ev)]
	//invoke
	evBinding.applyEvent(source, ev)
}
