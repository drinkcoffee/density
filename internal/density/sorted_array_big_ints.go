package density

import (
    "fmt"
    "math/big"
    "sort"
)

// SortedArrayBigInts is a sorted array of big integers.
type SortedArrayBigInts struct {
	ids[] big.Int
}


// NewSortedArrayBigInts is used to create a new instance.
func NewSortedArrayBigInts() *SortedArrayBigInts {
	var s = SortedArrayBigInts{}
	s.ids = []big.Int{ } 
	var _ sort.Interface = &s  // Enforce interface compliance
	return &s
}

// Add is used to add a single element to the array
func (s *SortedArrayBigInts) Add(newElement big.Int) {
	s.ids = append(s.ids, newElement)
	s.sort()
}

// AddMany is used to add a multiple elements to the array
func (s *SortedArrayBigInts) AddMany(newElements []big.Int) {
	for _, newElement := range newElements {
		s.ids = append(s.ids, newElement)
	}
	s.sort()
}


// Get one of the indices
// NOTE: There is no bounds checking on index!
func (s *SortedArrayBigInts) Get(index int) (big.Int) {
	return s.ids[index]
}

// PrintAll prints out all of the ids.
func (s *SortedArrayBigInts) PrintAll() {
	for i := 0; i < len(s.ids); i++ {
		fmt.Printf(" %5d, %#x\n", i, &s.ids[i])
	}
}


// Len is the number of elements in the collection.
// Needed to support sort.Interface
func (s *SortedArrayBigInts) Len() int {
	return len(s.ids)
}


// Less reports whether the element with
// index i should sort before the element with index j.
// Needed to support sort.Interface
func (s *SortedArrayBigInts)  Less(i, j int) bool {
	return s.ids[i].Cmp(&s.ids[j]) == -1
}
	
// Swap swaps the elements with indexes i and j.
// Needed to support sort.Interface
func (s *SortedArrayBigInts) Swap(i, j int) {
	temp := s.ids[i]
	s.ids[i] = s.ids[j]
	s.ids[j] = temp
}

// Sort orders the ids from lowest to highest
func (s *SortedArrayBigInts) sort() {
	sort.Sort(s)
}

