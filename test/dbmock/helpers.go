package dbmock

func JoinColumns(prefix string, columns []string) []string {
	joined := make([]string, len(columns))
	for i, col := range columns {
		joined[i] = prefix + "__" + col
	}

	return joined
}
