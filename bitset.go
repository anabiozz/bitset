package bitset

import "math/bits"

// BitSet ..
type BitSet struct {
	data []uint64
	size int
}

// Size ..
func (set *BitSet) Size() int {
	return set.size
}

// const ..
const (
	shift = 6
	mask  = 0x3f
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func index(n int) int {
	return n >> shift
}

func posVal(n int) uint64 {
	return 1 << uint(n&mask)
}

// NewBitSet ..
func NewBitSet(ns ...int) *BitSet {
	if len(ns) == 0 {
		return new(BitSet)
	}

	max := ns[0]
	for _, n := range ns {
		if n > max {
			max = n
		}
	}

	if max < 0 {
		return new(BitSet)
	}

	s := &BitSet{
		data: make([]uint64, index(max)+1),
	}

	for _, n := range ns {
		if n >= 0 {
			s.data[index(n)] |= posVal(n)
			s.size++
		}
	}

	return s
}

// Contains ..
func (set *BitSet) Contains(n int) bool {
	i := index(n)
	if i >= len(set.data) {
		return false
	}
	return set.data[i]&posVal(n) != 0
}

// Clear ..
func (set *BitSet) Clear(n int) *BitSet {
	if n < 0 {
		return set
	}

	i := index(n)
	if i >= len(set.data) {
		return set
	}

	if set.data[i]&posVal(n) != 0 {
		set.data[i] &^= posVal(n)
		set.size--
	}

	return set
}

func (set *BitSet) trim() {
	d := set.data
	n := len(d) - 1
	for n >= 0 && d[n] == 0 {
		n--
	}
	set.data = d[:n+1]
}

// Add ..
func (set *BitSet) Add(n int) *BitSet {
	if n < 0 {
		return set
	}

	i := index(n)
	if i >= len(set.data) {
		ndata := make([]uint64, i+1)
		copy(ndata, set.data)
		set.data = ndata
	}

	if set.data[i]&posVal(n) == 0 {
		set.data[i] |= posVal(n)
		// set.size++
	}

	return set
}

func (set *BitSet) computeSize() int {
	d := set.data
	n := 0
	for i, l := 0, len(d); i < l; i++ {
		if w := d[i]; w != 0 {
			n += bits.OnesCount64(w)
		}
	}

	return n
}

// Intersect ..
func (set *BitSet) Intersect(other *BitSet) *BitSet {
	minLen := min(len(set.data), len(other.data))

	intersectSet := &BitSet{
		data: make([]uint64, minLen),
	}

	for i := 0; i < minLen; i++ {
		intersectSet.data[i] = set.data[i] & other.data[i]
	}

	intersectSet.size = set.computeSize()

	return intersectSet
}

// Union ..
func (set *BitSet) Union(other *BitSet) *BitSet {
	var maxSet, minSet *BitSet

	if len(set.data) > len(other.data) {
		maxSet, minSet = set, other
	} else {
		maxSet, minSet = other, set
	}

	unionSet := &BitSet{
		data: make([]uint64, len(maxSet.data)),
	}

	minLen := len(minSet.data)
	copy(unionSet.data[minLen:], maxSet.data[minLen:])

	for i := 0; i < minLen; i++ {
		unionSet.data[i] = set.data[i] | other.data[i]
	}

	unionSet.size = unionSet.computeSize()
	return unionSet
}

// Difference ..
func (set *BitSet) Difference(other *BitSet) *BitSet {

	setLen := len(set.data)
	otherLen := len(other.data)

	differenceSet := &BitSet{
		data: make([]uint64, setLen),
	}

	minLen := setLen
	if setLen > otherLen {
		copy(differenceSet.data[otherLen:], set.data[otherLen:])
		minLen = otherLen
	}

	for i := 0; i < minLen; i++ {
		differenceSet.data[i] = set.data[i] &^ other.data[i]
	}

	differenceSet.size = differenceSet.computeSize()

	return differenceSet
}

// Visit ..
func (set *BitSet) Visit(do func(int) (skip bool)) (aborted bool) {
	d := set.data

	for i, len := 0, len(d); i < len; i++ {
		w := d[i]
		if w == 0 {
			continue
		}

		n := i << shift
		for w != 0 {
			b := bits.TrailingZeros64(w)
			if do(n + b) {
				return true
			}

			w &^= 1 << uint64(b)
		}
	}

	return false
}
