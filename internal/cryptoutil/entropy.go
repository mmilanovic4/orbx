package cryptoutil

import "math"

func ShannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0
	}

	var freq [256]int
	for _, b := range data {
		freq[b]++
	}

	n := float64(len(data))
	var entropy float64
	for _, count := range freq {
		if count == 0 {
			continue
		}
		p := float64(count) / n
		entropy -= p * math.Log2(p)
	}

	return entropy
}
