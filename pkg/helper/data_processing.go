package helper

func RemoveDuplicates(nums []int) []int {
	if nums == nil || len(nums) == 0 {
		return []int{}
	}

	unique := make(map[int]bool)
	result := []int{}

	for _, num := range nums {
		if !unique[num] {
			unique[num] = true
			result = append(result, num)
		}
	}

	return result
}
