package density

import (
    "fmt"
    "math/big"
    "sort"
)

// The maximum value of any value in the ring.
var max = new(big.Int)


// ValueRing represents a ring values.
type ValueRing struct {
	ids *SortedArrayBigInts
}


// NewValueRing is used to create a new instance.
func NewValueRing(wordSize int64) *ValueRing {
	max.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)
    fmt.Printf("Max: 0x%x\n", max)

	var v = ValueRing{}
	v.ids = NewSortedArrayBigInts()
	return &v
}

// GetIds returns the underlying id object
func (v *ValueRing) GetIds() (*SortedArrayBigInts) {
	return v.ids
}

// GetNumberRange returns the size of the number range
func (v *ValueRing) GetNumberRange() (*big.Int) {
	return max
}



// NumIdsInRange returns the number of ids in the range offset to offset + windowSize mod the number range
func (v *ValueRing) NumIdsInRange(offset, windowSize *big.Int) (int) {
	// Find the first index in the array where the id is larger than offset
	startIndex := sort.Search(v.ids.Len(), func(i int) bool { 
		temp := v.ids.Get(i)
		return temp.Cmp(offset) == 1})

	var endOffset = new(big.Int)
    endOffset.Add(offset, windowSize)
	endIndex := sort.Search(v.ids.Len(), func(i int) bool { 
		temp := v.ids.Get(i)
		return temp.Cmp(endOffset) == 1})
	extraGivenWrapAround := 0;
	if max.Cmp(endOffset) == -1 {
		endOffset.Sub(endOffset, max)
		extraGivenWrapAround = sort.Search(v.ids.Len(), func(i int) bool { 
			temp := v.ids.Get(i)
			return temp.Cmp(endOffset) == 1})
	}
	return endIndex - startIndex + extraGivenWrapAround
}


// GetRelativeIndex returns the index that is "jump" away from index "from"
func (v *ValueRing) GetRelativeIndex(from, jump int) (int) {
	result := (from + jump) % v.ids.Len()
    return result
}
