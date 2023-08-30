package slice

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func JoinInt64(elems []int64) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(elems[0], 10)
	}
	n := len(elems) - 1
	for i := 0; i < len(elems); i++ {
		n += len(strconv.FormatInt(elems[i], 10))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(strconv.FormatInt(elems[0], 10))
	for _, s := range elems[1:] {
		b.WriteString(",")
		b.WriteString(strconv.FormatInt(s, 10))
	}
	return b.String()
}
func JoinInt(elems []int) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(int64(elems[0]), 10)
	}
	n := len(elems) - 1
	for i := 0; i < len(elems); i++ {
		n += len(strconv.FormatInt(int64(elems[i]), 10))
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(strconv.FormatInt(int64(elems[0]), 10))
	for _, s := range elems[1:] {
		b.WriteString(",")
		b.WriteString(strconv.FormatInt(int64(s), 10))
	}
	return b.String()
}
func IsInSlice(ss []string, str string) bool {
	for _, s := range ss {
		if s == str {
			return true
		}
	}
	return false
}

// DivideBigSlice 对过长的切片进行分组.step 为多少个元素分为一组.
func DivideBigSlice(in []string, step int) [][]string {
	if len(in) == 0 {
		return nil
	}
	if len(in) < step || step <= 0 {
		return [][]string{in}
	}
	// count 分的组数
	count := len(in) / step
	if len(in)%step != 0 {
		count += 1
	}
	var result [][]string
	for i := 0; i < len(in); i += step {
		var inNew []string
		if i+step < len(in) {
			inNew = in[i : i+step]
		} else {
			inNew = in[i:]
		}
		result = append(result, inNew)
	}
	return result
}

// FilterStr 去重. 去掉空值.
func FilterStr(in []string) (out []string) {
	filterMap := make(map[string]struct{}, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == "" {
			continue
		}
		if _, ok := filterMap[in[i]]; ok {
			continue
		}
		out = append(out, in[i])
		filterMap[in[i]] = struct{}{}
	}
	return out
}

// FilterInt64 去重. 去0值.
func FilterInt64(in []int64) (out []int64) {
	filterMap := make(map[int64]struct{}, len(in))
	for i := 0; i < len(in); i++ {
		if in[i] == 0 {
			continue
		}
		if _, ok := filterMap[in[i]]; ok {
			continue
		}
		out = append(out, in[i])
		filterMap[in[i]] = struct{}{}
	}
	return out
}

// DivideBigSliceInt64 对过长的切片进行分组.
func DivideBigSliceInt64(in []int64, step int) [][]int64 {
	if len(in) < step || step <= 0 {
		return [][]int64{in}
	}
	// count 分的组数
	count := len(in) / step
	if len(in)%step != 0 {
		count += 1
	}
	var result [][]int64
	for i := 0; i < len(in); i += step {
		var inNew []int64
		if i+step < len(in) {
			inNew = in[i : i+step]
		} else {
			inNew = in[i:]
		}
		result = append(result, inNew)
	}
	return result
}

// union(并集)、intersect(交集)和 except(差集)

// Except 差集.
func Except(in []string, excepts ...string) []string {
	if len(in) == 0 || len(excepts) == 0 {
		return in
	}
	result := make([]string, 0, len(in))
	for _, v := range in {
		if IsInSlice(excepts, v) {
			continue
		}
		result = append(result, v)
	}
	return result
}

// Merge 并集。不合并merges的空函数.
func Merge(in []string, merges ...string) []string {
	if len(in) == 0 || len(merges) == 0 {
		return in
	}
	for _, v := range merges {
		if v == "" {
			continue
		}
		if IsInSlice(in, v) {
			continue
		}
		in = append(in, v)
	}
	return in
}

// GetSliceRandom 从数组中随即获取一个值.
func GetSliceRandom(arrays []string) string {
	if len(arrays) == 0 {
		return ""
	}
	return arrays[time.Now().UnixNano()%int64(len(arrays))]
}

// OrderStr 排序用的结构体.
type OrderStr struct {
	Value string
	Num   int64
}
type OrderStrs []OrderStr

func (strNumOrders OrderStrs) SortAndGet(isDesc bool) []string {
	if isDesc {
		sort.Slice(strNumOrders, func(i, j int) bool {
			if strNumOrders[i].Num > strNumOrders[j].Num {
				return true
			} else if strNumOrders[i].Num == strNumOrders[j].Num {
				return strNumOrders[i].Value < strNumOrders[j].Value
			} else {
				return false
			}
		})
	} else {
		sort.Slice(strNumOrders, func(i, j int) bool {
			if strNumOrders[i].Num < strNumOrders[j].Num {
				return true
			} else if strNumOrders[i].Num == strNumOrders[j].Num {
				return strNumOrders[i].Value < strNumOrders[j].Value
			} else {
				return false
			}
		})
	}
	items := make([]string, 0, len(strNumOrders))
	for _, order := range strNumOrders {
		items = append(items, order.Value)
	}
	return items
}

// OrderNum 排序用的结构体.
type OrderNum struct {
	Value int64
	Num   int64
}

type OrderNums []OrderNum

func (strNumOrders OrderNums) SortAndGet(isDesc bool) []int64 {
	if isDesc {
		sort.Slice(strNumOrders, func(i, j int) bool {
			if strNumOrders[i].Num > strNumOrders[j].Num {
				return true
			} else if strNumOrders[i].Num == strNumOrders[j].Num {
				return strNumOrders[i].Value < strNumOrders[j].Value
			} else {
				return false
			}
		})

	} else {
		sort.Slice(strNumOrders, func(i, j int) bool {
			if strNumOrders[i].Num < strNumOrders[j].Num {
				return true
			} else if strNumOrders[i].Num == strNumOrders[j].Num {
				return strNumOrders[i].Value < strNumOrders[j].Value
			} else {
				return false
			}
		})
	}
	items := make([]int64, 0, len(strNumOrders))
	for _, order := range strNumOrders {
		items = append(items, order.Value)
	}
	return items
}

// GenRandIV .
func GenRandIV() []byte {
	const length int = 16
	iv := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		iv[i] = byte(r.Int63n(255))
	}
	return iv
}
