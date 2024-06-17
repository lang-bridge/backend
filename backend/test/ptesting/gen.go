package ptesting

const alpha = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func (g *Gen) NextString(min, max int) string {
	n := min + g.r.IntN(max-min+1)
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = alpha[g.r.IntN(len(alpha))]
	}
	return string(buf)
}

func Array[T any](size int, g *Gen, f func(*Gen) T) []T {
	var res = make([]T, size)
	for i := 0; i < size; i++ {
		res[i] = f(g)
	}
	return res
}
