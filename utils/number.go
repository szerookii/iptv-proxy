package utils

// helpers functions for stb -> xtream otherwise it takes times or eat cpu and memory

func MergeNumbers(num1, num2 int32) int64 {
	return (int64(num1) << 32) | int64(num2)
}

func SplitNumbers(mergedNum int64) (int32, int32) {
	num1 := int32(mergedNum >> 32)
	num2 := int32(mergedNum & 0xFFFFFFFF)
	return num1, num2
}
