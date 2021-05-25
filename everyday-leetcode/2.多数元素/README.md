# 介绍

## 获取一个数组中只出现一次的值

### 题目说明

给定一个大小为 n 的数组，找到其中的多数元素。多数元素是指在数组中出现次数 大于 ⌊ n/2 ⌋ 的元素。

你可以假设数组是非空的，并且给定的数组总是存在多数元素。

示例 1：

输入：[3,2,3]
输出：3
示例 2：

输入：[2,2,1,1,1,2,2]
输出：2

### 解题思路

使用 go map 方式解决问题，通过遍历一次，再将 map 中存放元素的数值进行比较。

```go
func majorityElement(nums []int) int {
	value := make(map[int]int)

	// +1 的意思是为了能够匹配上数组中的元素
	count := len(nums) / 2 + 1

	for _, v := range nums {
		value[v]++
		if value[v] >= count {
			return v
		}
	}
	return 0
}
```

### 特殊解题思路

#### 摩尔投票法

1）当数组中的元素与假设的target不相等时，计数cnt减1，即模拟不同数字相互抵消；

2）假设数组中的元素与假设的target相等时，计数cnt加1；

3）当计数cnt等于0时，说明在当前遍历到的数组元素中，当前假设的target与其他数字相互抵消（个数相同），所以我们重新假设下一个遍历的数组元素为target,继续上面过程。

4）当遍历完数组后，target为所求数字。

```go
func majorityElement(nums []int) int {
	tn := nums[0]
	var cnt int
	for _, v := range nums {
		if tn == v {
			cnt += 1
		} else {
			cnt -= 1
		}
		if cnt == -1 {
			tn = v
			cnt = 0
		}
	}
	return tn
}
```