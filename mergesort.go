package main

import (
	"fmt"
	"math/rand"
	"time"
)

func makeRandomNumbers(num int) []int {

	var list []int

	list = make([]int, 0)

	for i := 0; i < num; i++ {

		list = append(list, int(rand.Float64()*100.0))
	}

	return list
}

func merge(left, right []int) []int, int {
	mergedArray := make([]int, 0, len(left)+len(right))
	inversions := 0

	for len(left) > 0 || len(right) > 0 {
		if len(left) == 0 {
			return append(mergedArray, right...)
		}
		if len(right) == 0 {
			return append(mergedArray, left...)
		}
		if left[0] < right[0] {
			mergedArray = append(mergedArray, left[0])
			left = left[1:]
		} else {
			mergedArray = append(mergedArray, right[0])
			right = right[1:]
		}
	}
	return mergedArray
}

func mergeSort(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	midpoint := len(list) / 2
	left := mergeSort(list[:midpoint])
	right := mergeSort(list[midpoint:])
	sorted := merge(left, right)

	return sorted
}

func main() {

	b := time.Now().UTC().UnixNano() % 10

	rand.Seed(b)

	start := time.Now()

	a := makeRandomNumbers(100)

	//fmt.Println(a)

	a = mergeSort(a)

	fmt.Println(a)

	delta := time.Now().Sub(start)

	fmt.Println("Execution time: ", delta)

	fmt.Println("seed used for random numbers: ", b)

}
