package filter

func indexOfString(collection []string, value string) int {
	for i, v := range collection {
		if v == value {
			return i
		}
	}

	return -1
}
