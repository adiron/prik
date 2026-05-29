package main

type BinArr8 uint8

func (b *BinArr8) flipBit(nth uint8) {
	*b ^= 1 << nth
}

func (b BinArr8) isBitFlipped(nth uint8) bool {
	return (b >> nth) & 0b01 == 1	
}
