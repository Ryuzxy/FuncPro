package fx


type Option[T any] struct {
    value T
    ok bool
}


func Some[T any](v T) Option[T] { return Option[T]{value: v, ok: true} }
func None[T any]() Option[T] { return Option[T]{ok: false} }


func (o Option[T]) IsSome() bool { return o.ok }
func (o Option[T]) IsNone() bool { return !o.ok }


func (o Option[T]) UnwrapOr(def T) T {
    if o.ok {
        return o.value
    }
    return def
}