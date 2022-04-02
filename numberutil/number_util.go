package numberutil

import "fmt"

// Factorial 階乗計算
func Factorial(n uint32) uint32 {
	if n == 0 {
		return 1
	}
	result := n
	for i := uint32(1); i < n-1; i++ {
		result *= n - i
	}
	return uint32(result)
}

// CombinationCount 組み合わせ数計算
func CombinationCount(n uint32, r uint32) (uint32, error) {
	permutation, err := PermutationCount(n, r)
	if err != nil {
		return 0, err
	}
	return permutation / Factorial(r), nil
}

// PermutationCount 順列数計算
func PermutationCount(n uint32, r uint32) (uint32, error) {
	if n == 0 || r == 0 {
		return 1, nil
	}
	if n < r {
		return 0, fmt.Errorf("r must be equal or less than n")
	}
	result := n
	for i := uint32(1); i < r; i++ {
		result *= n - i
	}
	return result, nil
}
