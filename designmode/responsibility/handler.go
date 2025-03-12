package responsibility

// Handler is the interface for a handler in the chain of responsibility pattern.
type Handler[I any, O any] interface {
	// Name returns the name of the handler.
	Name() string
	// Handle handles the request, returns true if the chain should stop.
	Handle(ctx *ChainCtx[I, O]) (stop bool, err error)
	// Rollback rolls back the request.
	Rollback(ctx *ChainCtx[I, O])
	// Pre returns the previous handler in the chain.
	Pre() Handler[I, O]
	// Next returns the next handler in the chain.
	Next() Handler[I, O]
	// SetPre sets the previous handler in the chain.
	SetPre(Handler[I, O])
	// SetNext sets the next handler in the chain.
	SetNext(Handler[I, O])
}

// HandlerBase implements the Handler interface except `Name()`.
type HandlerBase[I any, O any] struct {
	next Handler[I, O]
	pre  Handler[I, O]
}

func (h *HandlerBase[I, O]) Next() Handler[I, O] {
	return h.next
}

func (h *HandlerBase[I, O]) Pre() Handler[I, O] {
	return h.pre
}

func (h *HandlerBase[I, O]) SetNext(next Handler[I, O]) {
	h.next = next
}

func (h *HandlerBase[I, O]) SetPre(pre Handler[I, O]) {
	h.pre = pre
}

func (h *HandlerBase[I, O]) Rollback(ctx *ChainCtx[I, O]) {}
