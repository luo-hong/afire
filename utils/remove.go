package utils

import "sort"

func RemoveRepeatString(src []string) []string {
	end := removeRepeatSource(sort.StringSlice(src))
	return src[:end]
}

func RemoveRepeatInt(src []int) []int {
	end := removeRepeatSource(sort.IntSlice(src))
	return src[:end]
}

func RemoveRepeatInt64(src []int64) []int64 {
	end := removeRepeatSource(Int64Slice(src))
	return src[:end]
}

func removeRepeatSource(src sort.Interface) int {
	if src.Len() == 0 {
		return 0
	}
	sort.Sort(src)
	j := 0
	for i := 0; i < src.Len(); i++ {
		if src.Less(i, j) || src.Less(j, i) {
			// 两个不相等
			j++ // 替换位置后移
			if i != j {
				src.Swap(i, j)
			}
		}
	}
	return j + 1
}

type Int64Slice []int64

func (x Int64Slice) Len() int           { return len(x) }
func (x Int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
