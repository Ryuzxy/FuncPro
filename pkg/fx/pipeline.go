package fx

import (
    "context"
    "sync"
)

// PipelineStage represents a processing stage in pipeline
type PipelineStage[T any, R any] func(context.Context, T) Result[R]

// Pipeline represents a processing pipeline
type Pipeline[T any, R any] struct {
    stages []PipelineStage[interface{}, interface{}]
}

// NewPipeline creates new pipeline
func NewPipeline[T any, R any]() *Pipeline[T, R] {
    return &Pipeline[T, R]{
        stages: make([]PipelineStage[interface{}, interface{}], 0),
    }
}

// AddStage adds processing stage to pipeline
func (p *Pipeline[T, R]) AddStage[U any](stage PipelineStage[R, U]) *Pipeline[T, U] {
    p.stages = append(p.stages, func(ctx context.Context, input interface{}) Result[interface{}] {
        result := stage(ctx, input.(R))
        if result.IsErr() {
            return Err[interface{}](result.Unwrap().Error)
        }
        return Ok[interface{}](result.Unwrap().Value)
    })
    return &Pipeline[T, U]{stages: p.stages}
}

// Execute runs pipeline with input
func (p *Pipeline[T, R]) Execute(ctx context.Context, input T) Result[R] {
    var current interface{} = input
    
    for _, stage := range p.stages {
        result := stage(ctx, current)
        if result.IsErr() {
            return Err[R](result.Unwrap().Error)
        }
        current = result.Unwrap().Value
    }
    
    return Ok(current.(R))
}

// ParallelMap processes items in parallel with bounded concurrency
func ParallelMap[T any, R any](
    ctx context.Context,
    items []T,
    fn func(context.Context, T) Result[R],
    workers int,
) Result[[]R] {
    if workers <= 0 {
        workers = 1
    }
    
    type result struct {
        value R
        err   error
        index int
    }
    
    inputCh := make(chan T, len(items))
    resultCh := make(chan result, len(items))
    
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for item := range inputCh {
                select {
                case <-ctx.Done():
                    return
                default:
                    res, err := fn(ctx, item).Unwrap()
                    resultCh <- result{value: res, err: err}
                }
            }
        }()
    }
    
    // Send items to workers
    go func() {
        for _, item := range items {
            inputCh <- item
        }
        close(inputCh)
        
        wg.Wait()
        close(resultCh)
    }()
    
    // Collect results
    results := make([]R, 0, len(items))
    for res := range resultCh {
        if res.err != nil {
            return Err[[]R](res.err)
        }
        results = append(results, res.value)
    }
    
    return Ok(results)
}