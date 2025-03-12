package responsibility

// Builder is the builder for a chain of responsibility.
type Builder[I any, O any] struct {
	// inboundParam is the inbound parameter for the chain.
	inboundParam I
	// outboundFactory is the factory for the outbound response.
	outboundFactory OutboundFactory[O]
	// handlers is the handlers for the chain.
	handlers []Handler[I, O]
	// rollbackOnError is the flag to stop and rollback the chain on error.
	rollbackOnError bool
}

// Option is the option for the Builder.
type Option[I any, O any] func(*Builder[I, O])

// RollbackOnError sets the flag to rollback the chain on error.
func RollbackOnError[I any, O any](yes bool) Option[I, O] {
	return func(b *Builder[I, O]) {
		b.rollbackOnError = yes
	}
}

// OutboundFactory is the factory for the outbound response.
type OutboundFactory[O any] interface {
	// New creates a new outbound response.
	New() *O
}

func NewBuilder[I any, O any](inboundParam I, outboundFactory OutboundFactory[O], options ...Option[I, O]) *Builder[I, O] {
	b := &Builder[I, O]{
		inboundParam:    inboundParam,
		outboundFactory: outboundFactory,
	}
	for _, option := range options {
		option(b)
	}
	return b
}

// Link adds a handler to the chain.
func (b *Builder[I, O]) Link(handler Handler[I, O]) *Builder[I, O] {
	b.handlers = append(b.handlers, handler)
	return b
}

// Execute executes the chain.
func (b *Builder[I, O]) Execute() (*ChainCtx[I, O], error) {
	ctx := NewChainCtx(b.inboundParam, b.outboundFactory)
	for i, handler := range b.handlers {
		// First check if the handler can handle the request
		if !handler.CanHandle(ctx) {
			continue
		}

		// If the handler can handle the request, let it process
		stop, err := handler.Handle(ctx)
		if err != nil {
			ctx.Metadata[handler.Name()+"_error"] = err
			if b.rollbackOnError {
				// Rollback in reverse order from current handler
				for j := i - 1; j >= 0; j-- {
					if b.handlers[j].CanHandle(ctx) {
						b.handlers[j].Rollback(ctx)
					}
				}
				return nil, err
			}
		}
		if stop {
			return ctx, err
		}
	}
	return ctx, nil
}
