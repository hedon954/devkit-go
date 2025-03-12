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

// AddHandler adds a handler to the chain.
func (b *Builder[I, O]) AddHandler(handler Handler[I, O]) *Builder[I, O] {
	if len(b.handlers) > 0 {
		lastHandler := b.handlers[len(b.handlers)-1]
		handler.SetPre(lastHandler)
		lastHandler.SetNext(handler)
	}
	b.handlers = append(b.handlers, handler)
	return b
}

// Execute executes the chain.
func (b *Builder[I, O]) Execute() (*ChainCtx[I, O], error) {
	ctx := NewChainCtx(b.inboundParam, b.outboundFactory)
	for i, handler := range b.handlers {
		stop, err := handler.Handle(ctx)
		if err != nil {
			ctx.Metadata[handler.Name()+"_error"] = err
			if b.rollbackOnError {
				for j := i; j >= 0; j-- {
					b.handlers[j].Rollback(ctx)
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
