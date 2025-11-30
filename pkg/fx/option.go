package fx

// Option represents an optional value
type Option[T any] struct {
    value *T
}

// Some creates an Option with value
func Some[T any](value T) Option[T] {
    return Option[T]{value: &value}
}

// None creates an empty Option
func None[T any]() Option[T] {
    return Option[T]{value: nil}
}

// IsSome checks if Option contains value
func (o Option[T]) IsSome() bool {
    return o.value != nil
}

// IsNone checks if Option is empty
func (o Option[T]) IsNone() bool {
    return o.value == nil
}

// Unwrap returns the value or panics
func (o Option[T]) Unwrap() T {
    if o.value == nil {
        panic("unwrap on none option")
    }
    return *o.value
}

// UnwrapOr returns value or default
func (o Option[T]) UnwrapOr(defaultValue T) T {
    if o.value == nil {
        return defaultValue
    }
    return *o.value
}

// Map transforms the value if present
func (o Option[T]) Map[R any](fn func(T) R) Option[R] {
    if o.value == nil {
        return None[R]()
    }
    return Some(fn(*o.value))
}

// AndThen chains operations that return Option
func (o Option[T]) AndThen[R any](fn func(T) Option[R]) Option[R] {
    if o.value == nil {
        return None[R]()
    }
    return fn(*o.value)
}