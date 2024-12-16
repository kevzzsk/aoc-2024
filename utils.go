package utils

func AbsDiff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func Abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func RemoveElement(slice []int, index int) []int {
	if index < 0 || index >= len(slice) {
		return slice // return the original slice if the index is out of range
	}
	newSlice := make([]int, 0, len(slice)-1)
	newSlice = append(newSlice, slice[:index]...)
	return append(newSlice, slice[index+1:]...)
}

func CountBool(b ...bool) int {
	count := 0
	for _, v := range b {
		if v {
			count += 1
		}
	}
	return count
}
