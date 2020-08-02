package flux

import (
	"context"
	"fmt"

	"github.com/jjeffcaii/reactor-go"
)

const (
	BuffSizeXS = 32
	BuffSizeSM = 256
)

var _buffSize = BuffSizeSM

// InitBuffSize initialize the size of buff. (default=256)
func InitBuffSize(size int) {
	if size < 1 {
		panic(fmt.Sprintf("invalid flux buff size: %d", size))
	}
	_buffSize = size
}

type fluxCreate struct {
	source       func(context.Context, Sink)
	backpressure OverflowStrategy
}

func (m fluxCreate) SubscribeWith(ctx context.Context, s reactor.Subscriber) {
	var sink interface {
		reactor.Subscription
		Sink
	}
	switch m.backpressure {
	case OverflowDrop:
		// TODO: need implementation
		panic("implement me")
	case OverflowError:
		// TODO: need implementation
		panic("implement me")
	case OverflowIgnore:
		// TODO: need implementation
		panic("implement me")
	case OverflowLatest:
		// TODO: need implementation
		panic("implement me")
	default:
		sink = newBufferedSink(s, _buffSize)
	}
	s.OnSubscribe(sink)
	m.source(ctx, sink)
}

type CreateOption func(*fluxCreate)

func WithOverflowStrategy(o OverflowStrategy) CreateOption {
	return func(create *fluxCreate) {
		create.backpressure = o
	}
}

func newFluxCreate(c func(ctx context.Context, sink Sink), options ...CreateOption) *fluxCreate {
	fc := &fluxCreate{
		source:       c,
		backpressure: OverflowBuffer,
	}
	for i := range options {
		options[i](fc)
	}
	return fc
}
