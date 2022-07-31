package main

// import (
// 	"fmt"

// 	"gocv.io/x/gocv"
// )

// type ContourInfo struct {
// 	Index  int
// 	Next   int32
// 	Prev   int32
// 	Child  int32
// 	Parent int32
// 	Area   float64
// }
// type Container []ContourInfo

// func (pC *Container) Append(v gocv.Veci, a float64, i int) {
// 	c := ContourInfo{
// 		Index:  i,
// 		Next:   v[0],
// 		Prev:   v[1],
// 		Child:  v[2],
// 		Parent: v[3],
// 		Area:   a,
// 	}
// 	*pC = append(*pC, c)
// }

// func (C *ContourInfo) Show() {
// 	fmt.Printf(
// 		"%v [ n: %v, l: %v, c: %v, p: %v ] a: %v\n",
// 		C.Index,
// 		C.Next,
// 		C.Prev,
// 		C.Child,
// 		C.Parent,
// 		C.Area,
// 	)
// }

// func QuickSort(c Container) Container {
// 	if len(c) <= 1 {
// 		return c
// 	} else {
// 		l := float64(len(c))
// 		mid := math.Floor(l / 2)
// 		m := map[float64]int{
// 			c[0].a:        0,
// 			c[int(mid)].a: 1,
// 			c[len(c)].a:   2,
// 		}
// 		var tmp []float64
// 		for k, _ := range m {
// 			tmp = append(tmp, k)
// 		}
// 		pivot := c[m[tmp[1]]]

// 		fmt.Println("Pivot: ", pivot)

// 		return c
// 	}
// }

	// rand.Seed(time.Now().UnixMicro())
	// b := rand.Perm(10)
	// a := []float64{}
	// for i, v := range b{
	// 	a = append(a, float64(v))
	// 	fmt.Printf("%v ", i)
	// }
	// fmt.Println()
	// for _, v := range a{
	// 	fmt.Printf("%v ", v)
	// }
	// fmt.Println()

	// l := float64(len(a))
	// mid := int(math.Floor(l / 2))
	// fmt.Println("l/2: ", (l / 2))
	// fmt.Println("Mid: ", mid)

	// pMap := map[float64]int{
	// 	a[0]:        0,
	// 	a[mid]:      mid,
	// 	a[len(a)-1]: len(a) - 1,
	// }
	// fmt.Println("pMap: ", pMap)
	// var keys []float64
	// for k := range pMap {
	// 	keys = append(keys, k)
	// }
	// fmt.Println("keys: ", keys)
	// sort.Float64s(keys)
	// fmt.Println("sorted keys: ", keys)

	// pivot := a[pMap[keys[1]]]
	// fmt.Println("pivot: ", pivot)

	// less := []float64{}
	// equal := []float64{}
	// more := []float64{}

	// for i, v := range a{
	// 	if v < pivot{
	// 		less = append(less, a[i])
	// 	} else if v == pivot{
	// 		equal = append(equal, a[i])
	// 	} else if v > pivot{
	// 		more = append(more, a[i])
	// 	} else{
	// 		log.Fatal("Error Sorting Values")
	// 	}
	// }

	// fmt.Printf("less: %v\nequal: %v\nmore: %v\n", less, equal, more)