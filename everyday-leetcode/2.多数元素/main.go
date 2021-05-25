package main

import "fmt"

func majorityElement(nums []int) int {
	value := make(map[int]int)

	count := len(nums)/2 + 1

	for _, v := range nums {
		value[v]++
		if value[v] >= count {
			return v
		}
	}
	return 0
}

func main() {
	fmt.Println(majorityElement([]int{3, 1, 3}))
}
