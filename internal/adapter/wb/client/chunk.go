package client

func chunkInt64(in []int64, n int) [][]int64 {
	if n <= 0 {
		n = 1
	}
	if len(in) == 0 {
		return nil
	}
	out := make([][]int64, 0, (len(in)+n-1)/n)
	for i := 0; i < len(in); i += n {
		j := i + n
		if j > len(in) {
			j = len(in)
		}
		out = append(out, in[i:j])
	}
	return out
}
