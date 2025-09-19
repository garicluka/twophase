package twophase

func rotate_right[T any](arr []T, left int, right int) {
	temp := arr[right]

	for i := right; i > left; i-- {
		arr[i] = arr[i-1]
	}
	arr[left] = temp
}

func rotate_right_8[T any](arr *[8]T, left int, right int) {
	temp := arr[right]

	for i := right; i > left; i-- {
		arr[i] = arr[i-1]
	}
	arr[left] = temp
}

func rotate_right_12[T any](arr *[12]T, left int, right int) {
	temp := arr[right]

	for i := right; i > left; i-- {
		arr[i] = arr[i-1]
	}
	arr[left] = temp
}

func rotate_left[T any](arr []T, left int, right int) {
	temp := arr[left]

	for i := left; i < right; i++ {
		arr[i] = arr[i+1]
	}
	arr[right] = temp
}

func rotate_left_12[T any](arr *[12]T, left int, right int) {
	temp := arr[left]

	for i := left; i < right; i++ {
		arr[i] = arr[i+1]
	}
	arr[right] = temp
}

func c_nk(n int, k int) int {
	if n < k {
		return 0
	}

	if k > n/2 {
		k = n - k
	}

	s := 1
	i := n
	j := 1

	for i != n-k {
		s *= i
		s /= j
		i -= 1
		j += 1
	}

	return s
}
