package index

type BitSet []uint64

func NewBitSet(size int) BitSet {
	numWords := (size + 63) / 64
	return make(BitSet, numWords)
}

func (bs BitSet) Set(bitIndex int) {
	wordIndex := bitIndex / 64
	bitOffset := uint(bitIndex % 64)
	bs[wordIndex] |= 1 << bitOffset
}

func (bs BitSet) Reset(pos int) {
	wordIndex := pos / 64
	bitOffset := uint(pos % 64)
	bs[wordIndex] &^= 1 << bitOffset
}

func (bs BitSet) IsSet(bitIndex int) bool {
	wordIndex := bitIndex / 64
	bitOffset := uint(bitIndex % 64)
	return (bs[wordIndex] & (1 << bitOffset)) != 0
}
