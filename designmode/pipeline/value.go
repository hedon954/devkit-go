package pipeline

type Value[T any] interface {
	GetNext() Value[T]
	SetNext(Value[T])
	Invoke(T)
}

type ValueBase[T any] struct {
	next Value[T]
}

func (v *ValueBase[T]) GetNext() Value[T] {
	return v.next
}

func (v *ValueBase[T]) SetNext(next Value[T]) {
	v.next = next
}
