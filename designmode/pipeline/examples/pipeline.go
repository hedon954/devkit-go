package main

import (
	"fmt"
	"strings"

	"github.com/hedon954/devkit-go/designmode/pipeline"
)

func main() {
	s := "11 22 33"
	fmt.Println("before pipeline: s = " + s)

	p := pipeline.NewStandardPipeline(&TailValue{}, &FirstValue{}, &SecondValue{})
	p.Invoke(s)
}

type FirstValue struct {
	pipeline.ValueBase[string]
}

func (f *FirstValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "11", "first")
	fmt.Println("after first Value handled: s = " + s)
	f.GetNext().Invoke(s)
}

type SecondValue struct {
	pipeline.ValueBase[string]
}

func (sv *SecondValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "22", "second")
	fmt.Println("after second Value handled: s = " + s)
	sv.GetNext().Invoke(s)
}

type TailValue struct {
	pipeline.ValueBase[string]
}

func (t *TailValue) Invoke(s string) {
	s = strings.ReplaceAll(s, "33", "tail")
	fmt.Println("after tail Value handled: s = " + s)
}
