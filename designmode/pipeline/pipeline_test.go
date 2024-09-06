package pipeline

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FirstValue struct {
	ValueBase[string]
}

func (f *FirstValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "11", "first")
	fmt.Println("after first Value handled: s = " + s)
	f.GetNext().Invoke(s)
}

type SecondValue struct {
	ValueBase[string]
}

func (sv *SecondValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "22", "second")
	fmt.Println("after second Value handled: s = " + s)
	sv.GetNext().Invoke(s)
}

type TailValue struct {
	ValueBase[string]
}

func (t *TailValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "33", "tail")
	fmt.Println("after tail Value handled: s = " + s)
}

func TestStandardPipeline_ShouldWork(t *testing.T) {
	assert.NotPanics(t, func() {
		s := "11 22 33"
		p := NewStandardPipeline(&TailValue{}, &FirstValue{}, &SecondValue{}, &SecondValue{})
		p.Invoke(s)
		assert.NotNil(t, p.GetTail())
	})
}
