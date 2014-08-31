package sourcing

import (
	check "gopkg.in/check.v1"
	"reflect"
)

type NameChanged struct {
	Name string
}

type BiggestNameChanged struct {
	Name string
}

type user struct {
	Aggregate

	name            string
	biggestNameEver string

	//	numberOfTimesBiggestNameWasSet int
}

func (this *user) ChangeName(value string) {

	if this.name != value {

		this.On(&NameChanged{value})

		if len(value) > len(this.biggestNameEver) {

			this.On(&BiggestNameChanged{value})
		}
	}

}

func (this *user) onNameChanged(ev *NameChanged) {
	this.name = ev.Name
}

func (this *user) onBiggestNameChanged(ev *BiggestNameChanged) {
	this.biggestNameEver = ev.Name
}

func (t *FrameworkTests) TestIntegration(c *check.C) {

	b := t.f.Bind(
		reflect.TypeOf(&user{}),
		func() interface{} {
			return &user{}
		},
		func(s interface{}, a Aggregate) {
			s.(*user).Aggregate = a
		})

	b.On(reflect.TypeOf(&NameChanged{}),
		func(s interface{}, e interface{}) {
			s.(*user).onNameChanged(e.(*NameChanged))
		})

	b.On(reflect.TypeOf(&BiggestNameChanged{}),
		func(s interface{}, e interface{}) {
			s.(*user).onBiggestNameChanged(e.(*BiggestNameChanged))
		})

	evs1, _ := t.f.Update(reflect.TypeOf(&user{}), "myUserId",
		func(s interface{}) {
			s.(*user).ChangeName("john")
		})

	evs2, _ := t.f.Update(reflect.TypeOf(&user{}), "myUserId",
		func(s interface{}) {
			s.(*user).ChangeName("jon")
		})

	evs3, _ := t.f.Update(reflect.TypeOf(&user{}), "myUserId",
		func(s interface{}) {
			s.(*user).ChangeName("johny")
		})

	evs4, _ := t.f.Update(reflect.TypeOf(&user{}), "myUserId",
		func(s interface{}) {
			s.(*user).ChangeName("johny")
		})

	c.Assert(evs1, check.HasLen, 2)
	c.Assert(evs2, check.HasLen, 1)
	c.Assert(evs3, check.HasLen, 2)
	c.Assert(evs4, check.HasLen, 0)
}
