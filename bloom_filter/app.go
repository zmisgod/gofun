package bloom_filter

import "github.com/willf/bitset"

const DefaultSize = 2 << 24

var seeds = []uint{7, 11, 13, 31, 37, 61}

type BloomFilter struct {
	Set  *bitset.BitSet
	Func [6]SimpleHash
}

func NewBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	for i := 0; i < len(bf.Func); i++ {
		bf.Func[i] = SimpleHash{DefaultSize, seeds[i]}
	}
	bf.Set = bitset.New(DefaultSize)
	return bf
}

func (bf BloomFilter) Add(value string) {
	for _, f := range bf.Func {
		bf.Set.Set(f.hash(value))
	}
}

func (bf BloomFilter) Contains(value string) bool {
	if value == "" {
		return false
	}
	ret := true
	for _, f := range bf.Func {
		ret = ret && bf.Set.Test(f.hash(value))
	}
	return ret
}

type SimpleHash struct {
	Cap  uint
	Seed uint
}

func (s SimpleHash) hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.Seed + uint(value[i])
	}
	return (s.Cap - 1) & result
}
