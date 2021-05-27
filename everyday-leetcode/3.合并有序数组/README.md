// 本题的要求是，把nums1的前m项和nums2的前n项合并，放入nums1中。
func merge(nums1 []int, m int, nums2 []int, n int) {
    //把nums1复制到temp中
    temp := make([]int, m)
    copy(temp, nums1)

    t, j := 0, 0 //t为temp的索引，j为nums2的索引
    for i := 0; i < len(nums1); i++ {
        //当t大于temp的长度，那就是说temp全部放进去了nums1中，那剩下的就是放nums2剩余的值了
        if t >= len(temp) {
            nums1[i] = nums2[j]
            j++
            continue
        }
        //当j大于nums2的长度的时候，那就是说明nums2全部都放进去了nums1中，那剩下的就是放temp剩余的值了
        if j >= n {
            nums1[i] = temp[t]
            t++
            continue
        }
        //比较nums2与temp对应值的大小，小的那个就放进nums1中
        if nums2[j] <= temp[t] {
            nums1[i] = nums2[j]
            j++
        } else {
            nums1[i] = temp[t]
            t++
        }
    }
    fmt.Println(nums1)
}
