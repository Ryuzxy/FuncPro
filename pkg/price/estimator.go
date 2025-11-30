package price

// Pure function untuk rata-rata harga
func AveragePrice(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Higher-order function: bisa ganti metode estimasi
func Estimate(method func([]float64) float64, data []float64) float64 {
	return method(data)
}
