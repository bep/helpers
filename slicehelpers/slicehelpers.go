package slicehelpers

// Chunk splits the slice s into n number of chunks.
func Chunk[T any](s []T, n int) [][]T {
	if len(s) == 0 {
		return nil
	}
	var partitions [][]T
	sizeDefault := len(s) / n
	sizeBig := len(s) - sizeDefault*n
	size := sizeDefault + 1
	for i, idx := 0, 0; i < n; i++ {
		if i == sizeBig {
			size--
			if size == 0 {
				break
			}
		}
		partitions = append(partitions, s[idx:idx+size])
		idx += size
	}
	return partitions
}
