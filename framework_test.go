package sourcing

import (
	check "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type FrameworkTests struct {
	f Framework
}

var _ = check.Suite(&FrameworkTests{})

func (s *FrameworkTests) SetUpSuite(c *check.C) {

}

func (s *FrameworkTests) SetUpTest(c *check.C) {
	s.f = NewFramework(NewMemoryEventStore())
}

func (s *FrameworkTests) TearDownTest(c *check.C) {
	s.f = nil
}

func (s *FrameworkTests) TearDownSuite(c *check.C) {

}
