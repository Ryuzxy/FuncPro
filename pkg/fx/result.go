package fx

import "fmt"

// Result adalah tipe generic mirip Rust Result<T, E>.
type Result[T any] struct {
    value T
    err   error
}

// Ok membuat result sukses.
func Ok[T any](v T) Result[T] {
    return Result[T]{value: v}
}

// Err membuat result error.
func Err[T any](e error) Result[T] {
    var zero T
    return Result[T]{value: zero, err: e}
}

// IsErr mengecek apakah error.
func (r Result[T]) IsErr() bool { return r.err != nil }

// IsOk mengecek apakah tidak error.
func (r Result[T]) IsOk() bool { return r.err == nil }

// Unwrap mengembalikan value dan error.
func (r Result[T]) Unwrap() (T, error) { return r.value, r.err }

// FxMap memetakan value jika tidak error.
func FxMap[T, R any](r Result[T], fn func(T) R) Result[R] {
    if r.err != nil {
        return Err[R](r.err)
    }
    return Ok(fn(r.value))
}

// AndThen chaining seperti flatMap.
func AndThen[T, R any](r Result[T], fn func(T) Result[R]) Result[R] {
    if r.err != nil {
        return Err[R](r.err)
    }
    return fn(r.value)
}

// OrElse fallback value jika error.
func (r Result[T]) OrElse(v T) T {
    if r.err != nil {
        return v
    }
    return r.value
}

// String formatting.
func (r Result[T]) String() string {
    if r.err != nil {
        return fmt.Sprintf("Err(%v)", r.err)
    }
    return fmt.Sprintf("Ok(%v)", r.value)
}

// Try menjalankan fungsi dan membungkusnya ke Result.
func Try[T any](fn func() (T, error)) Result[T] {
    v, e := fn()
    if e != nil {
        return Err[T](e)
    }
    return Ok(v)
}

// Match pola branching ala Rust.
func Match[T, R any](r Result[T], onOk func(T) R, onErr func(error) R) R {
    if r.err != nil {
        return onErr(r.err)
    }
    return onOk(r.value)
}
