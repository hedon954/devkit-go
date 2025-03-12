package responsibility

// Handler is the interface for a handler in the chain of responsibility pattern.
type Handler[I any, O any] interface {
	// Name returns the name of the handler.
	Name() string
	// CanHandle checks if the handler can handle the request.
	CanHandle(ctx *ChainCtx[I, O]) bool
	// Handle handles the request, returns true if the chain should stop.
	Handle(ctx *ChainCtx[I, O]) (stop bool, err error)
	// Rollback rolls back the request.
	Rollback(ctx *ChainCtx[I, O])
}

// HandlerBase is the base implementation of the Handler interface.
type HandlerBase[I any, O any] struct{}

func (h *HandlerBase[I, O]) Rollback(ctx *ChainCtx[I, O]) {}
