package fx

import (
    "context"
    "sync"
)

// Ganti nama menjadi FxResult untuk menghindari konflik
type FxResult[T any] struct {
    Value T
    Error error
}

// FxOk creates a successful result
func FxOk[T any](value T) FxResult[T] {
    return FxResult[T]{Value: value}
}

// FxErr creates an error result
func FxErr[T any](err error) FxResult[T] {
    return FxResult[T]{Error: err}
}

// IsErr checks if the result is an error
func (r FxResult[T]) IsErr() bool {
    return r.Error != nil
}

// Unwrap returns the value and error (multiple return values)
func (r FxResult[T]) Unwrap() (T, error) {
    return r.Value, r.Error
}

// PipelineStage adalah function biasa, generic di function-nya saja
type PipelineStage[T any, R any] func(context.Context, T) FxResult[R]

// Pipeline bukan lagi generic per method; ia menyimpan tahap generik via closure
type Pipeline[I any, O any] struct {
    stages []func(context.Context, interface{}) FxResult[interface{}]
}

// NewPipeline membuat pipeline baru
func NewPipeline[I any]() *Pipeline[I, I] {
    return &Pipeline[I, I]{stages: []func(context.Context, interface{}) FxResult[interface{}]{}}
}

// AddStage ditulis sebagai function BIASA, bukan method generic
func AddStage[A any, B any](p *Pipeline[any, A], stage PipelineStage[A, B]) *Pipeline[any, B] {
    wrapped := func(ctx context.Context, v interface{}) FxResult[interface{}] {
        input := v.(A)
        res := stage(ctx, input)

        if res.IsErr() {
            return FxErr[interface{}](res.Error)
        }
        return FxOk[interface{}](res.Value)
    }

    newStages := append(p.stages, wrapped)
    return &Pipeline[any, B]{stages: newStages}
}

// Execute menjalankan pipeline
func Execute[I any, O any](p *Pipeline[I, O], ctx context.Context, input I) FxResult[O] {
    var current interface{} = input

    for _, st := range p.stages {
        r := st(ctx, current)
        if r.IsErr() {
            return FxErr[O](r.Error)
        }
        current = r.Value
    }

    return FxOk(current.(O))
}

// ParallelMap tetap aman
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

    // Send items to input channel
    go func() {
        for _, item := range items {
            inputCh <- item
        }
        close(inputCh)
    }()

    // Wait for all workers to finish and close result channel
    go func() {
        wg.Wait()
        close(resultCh)
    }()

    out := make([]R, 0, len(items))
    for res := range resultCh {
        if res.err != nil {
            return FxErr[[]R](res.err)
        }
        out = append(out, res.value)
    }

    return FxOk(out)
}