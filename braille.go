package main

import (
	"math"
)

const brailleStart = 0x2800

func makeBraille(b uint8) rune {
	var codePoint rune = brailleStart
	for i := 0; i < 8; i++ {
		if ((b >> i) & 0b1) == 1 {
			codePoint += int32(math.Pow(2, float64(i)))
		}
	}
	
	return codePoint
}
