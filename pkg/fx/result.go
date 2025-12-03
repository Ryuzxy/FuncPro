package fx

import "fmt"

// Result represents a computation that may fail
type Result[T any] struct {
    value T
    err   error
}

// Constructors
func Ok[T any](v T) Result[T] { 
    return Result[T]{value: v, err: nil} 
}

func Err[T any](e error) Result[T] { 
    var z T 
    return Result[T]{value: z, err: e} 
}

// Map applies function to successful value (sebagai function, bukan method)
func FxMap[T, R any](r Result[T], fn func(T) R) Result[R] {
    if r.err != nil { 
        return Err[R](r.err) 
    }
    return Ok(fn(r.value))
}

// AndThen chains operations that may fail (sebagai function, bukan method)
func AndThen[T, R any](r Result[T], fn func(T) Result[R]) Result[R] {
    if r.err != nil { 
        return Err[R](r.err) 
    }
    return fn(r.value)
}

// OrElse returns default value on error
func (r Result[T]) OrElse(defaultValue T) T {
    if r.err != nil {
        return defaultValue
    }
    return r.value
}

// Unwrap returns value and error
func (r Result[T]) Unwrap() (T, error) { 
    return r.value, r.err 
}

// IsOk checks if result is successful
func (r Result[T]) IsOk() bool { 
    return r.err == nil 
}

// IsErr checks if result is error
func (r Result[T]) IsErr() bool {
    return r.err != nil
}

// Match handles both success and failure cases
func (r Result[T]) Match(
    onSuccess func(T) interface{},
    onFailure func(error) interface{},
) interface{} {
    if r.err != nil {
        return onFailure(r.err)
    }
    return onSuccess(r.value)
}

// String implements Stringer interface
func (r Result[T]) String() string {
    if r.err != nil {
        return fmt.Sprintf("Err(%v)", r.err)
    }
    return fmt.Sprintf("Ok(%v)", r.value)
}

// Try executes a function that may fail and returns Result
func Try[T any](fn func() (T, error)) Result[T] {
    value, err := fn()
    if err != nil {
        return Err[T](err)
    }
    return Ok(value)
}