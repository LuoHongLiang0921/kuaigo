// @Description 


package krand

import "github.com/LuoHongLiang0921/kuaigo/pkg/util/kcast"

// GenRandomInt
//  @Description  生成范围在[start,end), 类型为int的随机数
//  @Param start
//  @Param end
//  @Return int
func GenRandomInt(start int, end int) int {
	// 生成随机数
	num := r.Intn((end - start)) + start

	return num
}

// GenRandomIntList
//  @Description  生成范围在[start,end), 类型为int的n个不重复随机数
//  @Param start
//  @Param end
//  @Param count
//  @Return []int
func GenRandomIntList(start int, end int, count int) []int {
	// 范围检查
	if end < start || (end-start) < count {
		return nil
	}

	// 存放结果的slice
	nums := make([]int, 0)

	for len(nums) < count {
		// 生成随机数
		num := r.Intn((end - start)) + start

		// 查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}

// GenRandomInt64
//  @Description  生成范围在[start,end), 类型为int64的随机数
//  @Param start
//  @Param end
//  @Return int64
func GenRandomInt64(start int64, end int64) int64 {
	// 生成随机数
	num := r.Int63n((end - start)) + start

	return num
}

// GenRandomInt64List
//  @Description  生成范围在[start,end), 类型为int64的n个不重复随机数
//  @Param start
//  @Param end
//  @Param count
//  @Return []int64
func GenRandomInt64List(start int64, end int64, count int) []int64 {
	// 范围检查
	if end < start || (end-start) < kcast.ToInt64(count) {
		return nil
	}

	// 存放结果的slice
	nums := make([]int64, 0)

	for len(nums) < count {
		// 生成随机数
		num := r.Int63n((end - start)) + start

		// 查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
