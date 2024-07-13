package main

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, v := range nums {
		other, ok := numMap[target-v]
		if ok {
			return []int{other, i}
		}
		numMap[v] = i
	}
	return []int{-1, -1}
}

func main() {}
