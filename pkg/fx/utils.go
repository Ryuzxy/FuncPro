package fx

// Map transforms a slice using pure function
func Map[T any, R any](arr []T, fn func(T) R) []R {
    res := make([]R, 0, len(arr))
    for _, v := range arr {
        res = append(res, fn(v))
    }
    return res
}

// Filter keeps elements that satisfy the predicate
func Filter[T any](arr []T, keep func(T) bool) []T {
    res := make([]T, 0, len(arr))
    for _, v := range arr {
        if keep(v) {
            res = append(res, v)
        }
    }
    return res
}

// Reduce folds slice into single value
func Reduce[T any, R any](arr []T, init R, fn func(R, T) R) R {
    acc := init
    for _, v := range arr {
        acc = fn(acc, v)
    }
    return acc
}

// Compose functions right-to-left
func Compose[T any, U any, V any](f func(U) V, g func(T) U) func(T) V {
    return func(x T) V {
        return f(g(x))
    }
}

// Pipe functions left-to-right  
func Pipe[T any](x T, fns ...func(T) T) T {
    result := x
    for _, fn := range fns {
        result = fn(result)
    }
    return result
}

// FlatMap maps and flattens the result
func FlatMap[T any, R any](arr []T, fn func(T) []R) []R {
    result := make([]R, 0)
    for _, v := range arr {
        result = append(result, fn(v)...)
    }
    return result
}

// Unique returns unique elements from slice
func Unique[T comparable](arr []T) []T {
    seen := make(map[T]bool)
    return Filter(arr, func(x T) bool {
        if !seen[x] {
            seen[x] = true
            return true
        }
        return false
    })
}

// GroupBy groups elements by key
func GroupBy[T any, K comparable](arr []T, keyFn func(T) K) map[K][]T {
    groups := make(map[K][]T)
    for _, item := range arr {
        key := keyFn(item)
        groups[key] = append(groups[key], item)
    }
    return groups
}