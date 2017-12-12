package array

func Index(arr []string, item string) int {
	if arr == nil || len(arr) == 0 {
		return -1
	}
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}

func Include(arr []string, item string) bool {
	return Index(arr, item) >= 0
}

func IsEqual(arr1 []string, arr2 []string) bool {
	if arr1 == nil && arr2 == nil {
		return true
	}
	if arr1 == nil || arr2 == nil || len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
