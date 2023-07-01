package util

func CompareArrayDifference(a, b []uint) []uint {
	diff := []uint{}
	bMap := make(map[uint]bool)
	for _, elem := range b {
		bMap[elem] = true
	}

	for _, elem := range a {
		if _, ok := bMap[elem]; !ok {
			diff = append(diff, elem)
		}
	}

	return diff
}
