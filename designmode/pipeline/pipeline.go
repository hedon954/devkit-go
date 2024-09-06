package pipeline

type Pipeline[T any] interface {
	GetHead() Value[T]
	GetTail() Value[T]
	SetTail(Value[T])
	AddValue(Value[T])
	Invoke(T)
}

type StandardPipeline[T any] struct {
	head Value[T]
	tail Value[T]
}

func NewStandardPipeline[T any](tail, head Value[T], others ...Value[T]) *StandardPipeline[T] {
	p := &StandardPipeline[T]{}
	p.SetTail(tail)
	p.AddValue(head)
	for _, v := range others {
		p.AddValue(v)
	}
	return p
}

func (s *StandardPipeline[T]) Invoke(t T) {
	s.GetHead().Invoke(t)
}

func (s *StandardPipeline[T]) GetHead() Value[T] {
	return s.head
}

func (s *StandardPipeline[T]) GetTail() Value[T] {
	return s.tail
}

func (s *StandardPipeline[T]) SetTail(value Value[T]) {
	s.tail = value
}

func (s *StandardPipeline[T]) AddValue(value Value[T]) {
	if s.head == nil {
		s.head = value
		value.SetNext(s.tail)
	} else {
		cur := s.head
		for cur != nil {
			if cur.GetNext() == s.tail {
				cur.SetNext(value)
				value.SetNext(s.tail)
				break
			}
			cur = cur.GetNext()
		}
	}
}
