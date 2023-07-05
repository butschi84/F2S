package operator

func stringArrayContains(target string, arr []string) bool {
	found := false
	for _, str := range arr {
		if str == target {
			found = true
			break
		}
	}
	return found
}
