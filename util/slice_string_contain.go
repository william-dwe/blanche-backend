package util

func IsSliceContainString(strSlice []string, str string) bool {
	for _, s := range strSlice {
		if s == str {
			return true
		}
	}
	return false
}
