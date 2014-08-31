package sourcing

import (
	"reflect"
)

type Framework interface {
	Bind(sourceType reflect.Type, create CreateHandler, init InitHandler) Binding
	BindNamed(sourceType reflect.Type, sourceName string, create CreateHandler, init InitHandler) Binding

	//Load(id string, entity interface{}) (Aggregate, error)
	//Save(aggregate Aggregate) ([]*Event, error)

	Update(sourceType reflect.Type, id string, update UpdateHandler) ([]*Event, error)
}

type UpdateHandler func(entity interface{})

type framework struct {
	store    EventStore
	bindings map[reflect.Type]Binding
}

func NewFramework(store EventStore) Framework {
	return &framework{
		store:    store,
		bindings: make(map[reflect.Type]Binding),
	}
}

func (this *framework) Bind(sourceType reflect.Type, create CreateHandler, init InitHandler) Binding {

	return this.BindNamed(sourceType, sourceType.Name(), create, init)
}

func (this *framework) BindNamed(sourceType reflect.Type, sourceName string, create CreateHandler, init InitHandler) Binding {

	binding := newBinding(sourceType, sourceName, create, init)

	this.bindings[sourceType] = binding
	return binding
}

func (this *framework) load(id string, entity interface{}) (Aggregate, error) {

	events, err := this.store.Load(id)

	if err != nil {
		return nil, err
	}

	binder := this.bindings[reflect.TypeOf(entity)]

	aggregate := loadAggregate(binder, id, entity, events)

	return aggregate, nil
}

func (this *framework) save(aggregate Aggregate) ([]*Event, error) {

	//on conflicts due to original version being different than whats stored
	//load the new events only, apply on aggregate and try again until succeed

	changes := aggregate.Changes()

	err := this.store.Append(aggregate.Id(), aggregate.Version(), changes)

	if err == nil {
		return changes, nil
	}

	return nil, err
}

func (this *framework) Update(sourceType reflect.Type, id string, update UpdateHandler) ([]*Event, error) {

	binding := this.bindings[sourceType]

	entity := binding.Create()
	aggregate, err := this.load(id, entity)
	binding.Init(entity, aggregate)

	if err != nil {
		return nil, err
	}

	update(entity)

	changes, err2 := this.save(aggregate)

	return changes, err2
}
