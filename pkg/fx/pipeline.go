package fx

import (
    "context"
    "sync"
)

type FxResult[T any] struct {
    Value T
    Error error
}

func FxOk[T any](value T) FxResult[T]  { return FxResult[T]{Value: value} }
func FxErr[T any](err error) FxResult[T] { return FxResult[T]{Error: err} }
func (r FxResult[T]) IsErr() bool       { return r.Error != nil }
func (r FxResult[T]) Unwrap() (T, error) { return r.Value, r.Error }

type PipelineStage[T any, R any] func(context.Context, T) FxResult[R]

type Pipeline[I any, O any] struct {
    stages []func(context.Context, interface{}) FxResult[interface{}]
}

func NewPipeline[I any]() *Pipeline[I, I] {
    return &Pipeline[I, I]{stages: []func(context.Context, interface{}) FxResult[interface{}]{}}
}

func AddStage[A any, B any](p *Pipeline[any, A], stage PipelineStage[A, B]) *Pipeline[any, B] {
    wrapper := func(ctx context.Context, v interface{}) FxResult[interface{}] {
        input := v.(A)
        res := stage(ctx, input)
        if res.IsErr() {
            return FxErr[interface{}](res.Error)
        }
        return FxOk[interface{}](res.Value)
    }

    return &Pipeline[any, B]{stages: append(p.stages, wrapper)}
}

func Execute[I any, O any](p *Pipeline[I, O], ctx context.Context, input I) FxResult[O] {
    var cur interface{} = input

    for _, st := range p.stages {
        r := st(ctx, cur)
        if r.IsErr() {
            return FxErr[O](r.Error)
        }
        cur = r.Value
    }

    return FxOk(cur.(O))
}

func ParallelMap[T any, R any](
    ctx context.Context,
    items []T,
    fn func(context.Context, T) FxResult[R],
    workers int,
) FxResult[[]R] {

    if workers <= 0 {
        workers = 1
    }

    type result struct {
        value R
        err   error
    }

    inputCh := make(chan T, len(items))
    resultCh := make(chan result, len(items))

    var wg sync.WaitGroup

    // workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range inputCh {
                select {
                case <-ctx.Done():
                    return
                default:
                    res := fn(ctx, item)
                    if res.IsErr() {
                        resultCh <- result{err: res.Error}
                    } else {
                        resultCh <- result{value: res.Value}
                    }
                }
            }
        }()
    }

    // send items
    go func() {
        for _, item := range items {
            inputCh <- item
        }
        close(inputCh)
    }()

    // wait for workers then close resultCh
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    output := make([]R, 0, len(items))
    for res := range resultCh {
        if res.err != nil {
            return FxErr[[]R](res.err)
        }
        output = append(output, res.value)
    }

    return FxOk(output)
}
