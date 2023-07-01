package util

func FindIdxArrString(key string, variantType []string) int {
	for idx, variantType := range variantType {
		if variantType == key {
			return idx
		}
	}
	return -1
}
