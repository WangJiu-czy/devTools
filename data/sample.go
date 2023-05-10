package data

import "math/rand"

func Samples[T any](collection []T, count int) []T {
	size := len(collection)

	ts := append([]T{}, collection...)
	results := []T{}

	for i := 0; i < size && i < count; i++ {
		copyLength := size - i
		index := rand.Intn(size - i)
		results = append(results, ts[index])

		//ts[index]的值已经用过了,所以将最后一个值移动过来,覆盖掉用过的值
		ts[index] = ts[copyLength-1]
		//缩小取样切片的长度,删除刚刚移动的最后一个值
		ts = ts[:copyLength-1]
	}
	return results
}
