package density

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "sort"
)

var max = new(big.Int)


// GatewayIds is a sortable array of big integers.
type GatewayIds struct {
	ids[] big.Int
	sorted bool
}


// NewGatewayIds is used to create a new instance.
func NewGatewayIds(wordSize int64) *GatewayIds {
	max.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)
    fmt.Printf("Max: 0x%x\n", max)

	var gatewayIds = GatewayIds{}
	gatewayIds.ids = []big.Int{ } 
	gatewayIds.sorted = true

	var _ sort.Interface = &gatewayIds  // Enforce interface compliance
	return &gatewayIds
}

// AddRandom is used to generate a random id in the range 0 to 256**wordSize - 1
func (g *GatewayIds) AddRandom() {
	n, _ := rand.Int(rand.Reader, max)
	g.ids = append(g.ids, *n)
//	fmt.Printf(" New randomly generated Gateway Id: 0x%x\n", &g.ids[len(g.ids) - 1])
	g.sorted = false
}

// Get one of the indices
// NOTE: There is no bounds checking on index!
func (g *GatewayIds) Get(index int) (big.Int) {
	return g.ids[index]
}

// PrintAll prints out all of the ids.
func (g *GatewayIds) PrintAll() {
	for i := 0; i < len(g.ids); i++ {
		fmt.Printf(" Gateway Id: 0x%x\n", &g.ids[i])
	}
}

// NumIdsInRange returns the number of ids in the range offset to offset + windowASize mod the number range
func (g *GatewayIds) NumIdsInRange(offset, windowSize *big.Int) (int) {
	// Find the first index in the array where the id is larger than offset
	startIndex := sort.Search(g.Len(), func(i int) bool { return g.ids[i].Cmp(offset) == 1})

	var endOffset = new(big.Int)
    endOffset.Add(offset, windowSize)
	endIndex := sort.Search(g.Len(), func(i int) bool { return g.ids[i].Cmp(endOffset) == 1})
	extraGivenWrapAround := 0;
	if max.Cmp(endOffset) == -1 {
		// fmt.Printf(" max: 0x%x\n", max)
		// fmt.Printf(" endOffset: 0x%x\n", endOffset)
		endOffset.Sub(endOffset, max)
		extraGivenWrapAround = sort.Search(g.Len(), func(i int) bool { return g.ids[i].Cmp(endOffset) == 1})
	}
	return endIndex - startIndex + extraGivenWrapAround
}


// Sort orders the ids from lowest to highest
func (g *GatewayIds) Sort() {
	sort.Sort(g)
	g.sorted = true
}

func (g *GatewayIds) isSorted() (bool) {
	return g.sorted
}

// Len is the number of elements in the collection.
// Needed to support sort.Interface
func (g *GatewayIds) Len() int {
	return len(g.ids)
}


// Less reports whether the element with
// index i should sort before the element with index j.
// Needed to support sort.Interface
func (g *GatewayIds)  Less(i, j int) bool {
	return g.ids[i].Cmp(&g.ids[j]) == -1
}
	
// Swap swaps the elements with indexes i and j.
// Needed to support sort.Interface
func (g *GatewayIds) Swap(i, j int) {
	temp := g.ids[i]
	g.ids[i] = g.ids[j]
	g.ids[j] = temp
}
