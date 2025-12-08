package fx


func Map[T any, R any](arr []T, fn func(T) R) []R {
    out := make([]R, 0, len(arr))
    for _, v := range arr {
        out = append(out, fn(v))
    }
    return out
}


func Filter[T any](arr []T, keep func(T) bool) []T {
    out := make([]T, 0, len(arr))
    for _, v := range arr {
        if keep(v) {
            out = append(out, v)
        }
    }
    return out
}


func GroupBy[T any, K comparable](arr []T, keyFn func(T) K) map[K][]T {
    m := make(map[K][]T)
    for _, item := range arr {
        k := keyFn(item)
        m[k] = append(m[k], item)
    }
    return m
}