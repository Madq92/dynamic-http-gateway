package utils

// Contains check whether `slice` contain `target`
func Contains(slice []string, target string) bool {
	return IndexOf(slice, target) >= 0
}

// IndexOf for `target` index in `slice`, return -1 if `target` not found
func IndexOf(slice []string, target string) int {
	for i, item := range slice {
		if item == target {
			return i
		}
	}
	return -1
}
