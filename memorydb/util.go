package memorydb

// insert element el in arr and returns array and insertion postion.
// if multiple el exists, then select that last occurance

func sortedinsert(arr []int, el int) ([]int, int) {
	if len(arr) == 0 {
		arr = append(arr, el)
		return arr, 0
	}

	// binary search based on score
	low, high := 0, len(arr)
	for low < high {
		mid := (low + high) / 2
		// get the last occurance of element
		if arr[mid] >= el {
			low = mid + 1
		} else {
			high = mid
		}
	}

	// get the last element as it will be overridden by last-1 element
	if high == len(arr) {
		arr = append(arr, el)
	} else {
		last := arr[len(arr)-1]
		// increment position by 1 of elements having position >= ans
		copy(arr[high+1:], arr[high:])
		// append the last element, causing underlying array to increase if necessary
		arr = append(arr, last)
		// set the element
		arr[high] = el
	}

	return arr, high
}
