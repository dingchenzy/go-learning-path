# 介绍

## 获取一个数组中只出现一次的值

### 题目说明

给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。

说明：

你的算法应该具有线性时间复杂度。 你可以不使用额外空间来实现吗？

示例 1:

输入: [2,2,1]
输出: 1
示例 2:

输入: [4,1,2,1,2]
输出: 4

### 第一次解题思路

这简单，直接使用 for 循环遍历，并赋值给 map，让 map 自己进行 ++ 迭代，之后进行比较 map 的 value 值，如果为 1 那么肯定就是只出现一次的值啦。

```go
func singleNumber(nums []int) int {
    mp := make(map[int]int)
    for _, v := range nums {
        mp[v]++
    }
    for i, v := range mp {
        if v == 1 {
            return i
        }
    }
    return 0
}
```

提交后的执行效率却不尽人意，于是开始看题解，原来要使用异或的方式。

### 第二次解题思路

通过异或可以让相同的数等于 0，而不同的数异或得出的内容如果与之前相同的数相异或得出的结果依然是不同的那个值。

a ^ b ^ a = ( a ^ a ) ^ b = 0 ^ b = b

```go
func singleNumber(nums []int) int {
    var value int
    for _, v := range nums{
        value ^= v
    }
    return value
}
```