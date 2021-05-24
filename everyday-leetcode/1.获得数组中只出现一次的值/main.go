package main

import "fmt"

var (
	nums = []int{2, 1, 2}
)

func singleNumber(nums []int) int {
	var value int
	for _, v := range nums {
		value ^= v
	}
	return value
}

// func singleNumber(nums []int) int {
// mp := make(map[int]int)
// for _, v := range nums {
// mp[v]++
// }
// for i, v := range mp {
// if v == 1 {
// return i
// }
// }
// return 0
// }

func main() {
	fmt.Println(singleNumber(nums))
}
