package flux

import (
	"context"
	"fmt"
	"time"
)

var empty = just(nil)

func Error(e error) Flux {
	return wrap(newFluxError(e))
}

func Range(start, count int) Flux {
	if count < 0 {
		panic(fmt.Errorf("count >= 0 but is was %d", count))
	}
	if count == 0 {
		return empty
	}
	return wrap(newFluxRange(start, start+count))
}

func Empty() Flux {
	return empty
}

func Just(values ...Any) Flux {
	if len(values) < 1 {
		return empty
	}
	return just(values)
}

func Create(c func(ctx context.Context, sink Sink), options ...CreateOption) Flux {
	return wrap(newFluxCreate(c, options...))
}

func Interval(period time.Duration) Flux {
	return wrap(newFluxInterval(period, 0, nil))
}

func NewUnicastProcessor() Processor {
	return wrap(newUnicastProcessor(BuffSizeSM))
}

func just(values []Any) Flux {
	return wrap(newSliceFlux(values))
}
