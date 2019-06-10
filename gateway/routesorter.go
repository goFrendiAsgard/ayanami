package gateway

// RouteSorter a sorter for http routes
type RouteSorter []string

func (r RouteSorter) Len() int {
	return len(r)
}

func (r RouteSorter) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RouteSorter) Less(i, j int) bool {
	return len(r[i]) < len(r[j])
}
