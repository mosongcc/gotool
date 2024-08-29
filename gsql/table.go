package gsql

type SelectBuild struct {
	items   []string
	where   string
	groupBy string
	orderBy string
	limit   uint64
	offset  uint64
}

func And(k, j, v string) string {
	return "and" + k + j + v
}
