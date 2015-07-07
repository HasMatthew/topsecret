package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	//generate a random sequence with given length
	unsort := randomSequence(9)
	fmt.Println(unsort)

	//define the sorting alogrithm and do sorts

	//insertionSort(unsort)

	//slectionSort(unsort)

	//unsort = mergeSort(unsort)
	quickSort(unsort, 0, len(unsort)-1)

	//heapSort(unsort)

	fmt.Println(unsort)
}

func insertionSort(seq []int) {
	if len(seq) > 1 {
		for i := 1; i < len(seq); i++ {
			each := seq[i]
		outerloop:
			for j := 0; j < i; j++ {
				if each < seq[j] {
					for z := i; z > j; z-- {
						seq[z] = seq[z-1]
					}
					seq[j] = each
					break outerloop
				}
			}
		}
	}
}

func slectionSort(seq []int) {
	if len(seq) > 1 {
		for i := 0; i < len(seq); i++ {
			min := seq[i]
			index := i
			change := false
			for j := i + 1; j < len(seq); j++ {
				if min > seq[j] {
					min = seq[j]
					index = j
					change = true
				}
			}
			if change {
				temp := seq[i]
				seq[i] = seq[index]
				seq[index] = temp
			}
		}
	}
}

func mergeSort(seq []int) []int {
	if len(seq) > 1 {
		mid := len(seq) / 2
		left := seq[:mid]
		right := seq[mid:]
		left = mergeSort(left)
		right = mergeSort(right)
		return merge(left, right)
	} else {
		return seq
	}
}

func merge(left, right []int) []int {
	length := len(left) + len(right)
	a := make([]int, length)
	var i, j int
	for k := 0; k < length; k++ {
		if i < len(left) && j < len(right) && left[i] <= right[j] {
			a[k] = left[i]
			i++
		} else if i < len(left) && j < len(right) && left[i] > right[j] {
			a[k] = right[j]
			j++
		} else if i >= len(left) {
			a[k] = right[j]
			j++
		} else {
			a[k] = left[i]
			i++
		}
	}
	return a

}

func quickSort(seq []int, start, end int) {
	if start < end {
		//pivot index
		pivot := partition(seq, start, end)
		quickSort(seq, start, pivot-1)
		quickSort(seq, pivot+1, end)
	}
}

func partition(seq []int, start, end int) int {
	value := seq[end]
	i := start
	j := end
	for i < j {
		for i < j && seq[i] < value {
			i++
		}

		for i < j && seq[j] >= value {
			j--
		}

		if i < j {
			temp := seq[i]
			seq[i] = seq[j]
			seq[j] = temp
		}
	}
	seq[end] = seq[j]
	seq[j] = value
	return j

}

func heapSort(seq []int) {
	for i := (len(seq) - 1) / 2; i >= 0; i-- {
		binaryMinHeap(i, seq)
	}
}

func binaryMinHeap(hole int, seq []int) {
	var index int = seq[hole]
outer:
	for (2*hole + 1) < len(seq) {
		pre := hole
		left := 2*hole + 1
		right := 2*hole + 2
		if right < len(seq) && seq[right] < seq[left] && seq[right] < index {
			hole = right
		} else if seq[left] < index {
			hole = left
		} else {
			break outer
		}

		seq[pre] = seq[hole]
	}

	fmt.Println(hole)
	seq[hole] = index
}

func randomSequence(length int) []int {
	a := make([]int, length)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		a[i] = rand.Int() % 100
	}
	return a
}
