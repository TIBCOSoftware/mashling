// +build noasm

package gorgonia

import "math/bits"

func divmod(a, b int) (q, r int) {
	return a / b, a % b
}

func popcnt(a uint64) int {
	return bits.OnesCount64(a)
}

func clz(a uint64) int {
	return bits.LeadingZeros64(a)
}

func ctz(a uint64) int {
	return bits.TrailingZeros64(a)
}
