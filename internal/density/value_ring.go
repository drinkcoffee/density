package density

import (
    "fmt"
    "math/big"
    "sort"
)



// ValueRing represents a ring values.
type ValueRing struct {
	ids *SortedArrayBigInts
	
    max *big.Int            // The maximum value of any value in the ring.
   
    wordSize int64          // size in bytes
    dhtSpan int64
}


// NewValueRing is used to create a new instance.
func NewValueRing(wordSize, dhtSpan int64, initialArrayOfIndetifiers *SortedArrayBigInts) *ValueRing {
	var v = ValueRing{}
    v.ids = initialArrayOfIndetifiers
    v.wordSize = wordSize
    v.dhtSpan = dhtSpan
    v.max = big.NewInt(0)
	v.max.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)
	return &v
}

// GetIds returns the underlying id object
func (v *ValueRing) GetIds() (*SortedArrayBigInts) {
	return v.ids
}

// GetNumberRange returns the size of the number range
func (v *ValueRing) GetNumberRange() (*big.Int) {
	return v.max
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
	if v.max.Cmp(endOffset) == -1 {
		endOffset.Sub(endOffset, v.max)
		extraGivenWrapAround = sort.Search(v.ids.Len(), func(i int) bool { 
			temp := v.ids.Get(i)
			return temp.Cmp(endOffset) == 1})
	}
	return endIndex - startIndex + extraGivenWrapAround
}


// FindNextHighest finds the value that is the next highest in the ring
func (v *ValueRing) FindNextHighest(offset *big.Int) (int, big.Int) {
	foundIndex := sort.Search(v.ids.Len(), func(i int) bool { 
		temp := v.ids.Get(i)
		return temp.Cmp(offset) == 1})
	if foundIndex == v.ids.Len() {
		foundIndex = 0
	}
	return foundIndex, v.ids.Get(foundIndex)
}

// GetRelativeIndex returns the index that is "jump" away from index "from"
func (v *ValueRing) GetRelativeIndex(from, jump int) (int) {
	result := (from + jump) % v.ids.Len()
    return result
}


// FindBestLocationForValue finds the ideal location for a new value in the ring
// Ideal is defined as the middle of the largest gap between two identifiers in 
// the largest region of low dentity.
func (v *ValueRing) FindBestLocationForValue() (big.Int) {
    fmt.Println("Scan the number range to see how many identifiers cover each region.")
    fmt.Println("Determine the indices that correspond to regions that have the lowest number identifiers.")
    numRange := v.GetNumberRange()
    fmt.Printf(" Number Range: %#x\n", numRange)
    
    numNodes := v.ids.Len()
    var windowSize = new(big.Int)
    windowSize.Mul(numRange, big.NewInt(v.dhtSpan)).Div(windowSize, big.NewInt(int64(numNodes)))
    fmt.Printf(" Window Size:  %#x\n", windowSize)

    numSteps := numNodes * 2
    var stepSize = new(big.Int)
    stepSize.Div(numRange, big.NewInt(int64(numSteps)))
    fmt.Printf(" Step Size:    %#x\n", stepSize)

    fmt.Println(" Determine the weightings of each window")
    var indicesWithLowestNumIdsInRange []int = nil
    lowestNumIdsInRange := numNodes
    var offset = *big.NewInt(0)
    index := 0
    fmt.Println(" index, num identifiers in range")
    for offset.Cmp(numRange) == -1 {
        numIdsInRange := v.NumIdsInRange(&offset, windowSize)
        // Print out comma separated values of index of the scan and the number of ids in the range.
        fmt.Printf(" %5d, %d\n", index, numIdsInRange)

        if (numIdsInRange == int(lowestNumIdsInRange)) {
            indicesWithLowestNumIdsInRange = append(indicesWithLowestNumIdsInRange, index)
        }

        if (numIdsInRange < int(lowestNumIdsInRange)) {
            lowestNumIdsInRange = numIdsInRange
            indicesWithLowestNumIdsInRange = nil
            indicesWithLowestNumIdsInRange = append(indicesWithLowestNumIdsInRange, index)
        }

        index++
        offset.Add(&offset, stepSize)
    }

    fmt.Printf("List of windows/regions start offsets that have the lowest number of idenfiers (%d)\n", lowestNumIdsInRange)
    for _, s := range indicesWithLowestNumIdsInRange {
        fmt.Println(s)
    }
    fmt.Println()

    fmt.Println("Find the largest low density region (longest contiguous set of indices)")
    startOfLongestContiguousRun := 0
    lenOfLongestContiguousRun := 1
    lenOfContiguousRun := 1
    inContiguousRun := true
    numIndicesWithLowestNumIdsInRange := len(indicesWithLowestNumIdsInRange)
    // Go through the array twice to ensure we handle the wrap-around
    for i := 1; i < numIndicesWithLowestNumIdsInRange * 2; i++ {
        idx := i % numIndicesWithLowestNumIdsInRange
        idxLessOne := (i - 1) % numIndicesWithLowestNumIdsInRange
        if (indicesWithLowestNumIdsInRange[idx] - indicesWithLowestNumIdsInRange[idxLessOne] == 1) ||
           (indicesWithLowestNumIdsInRange[idxLessOne] == int(numSteps)-1 && indicesWithLowestNumIdsInRange[idx] == 0) {
            // Is contiguous
            if (inContiguousRun) {
                lenOfContiguousRun++
            } else {
                lenOfContiguousRun = 1
                inContiguousRun = true
            }
            if (lenOfLongestContiguousRun < lenOfContiguousRun) {
                lenOfLongestContiguousRun = lenOfContiguousRun
                startOfLongestContiguousRun = (idx + 1 - lenOfContiguousRun)
                if startOfLongestContiguousRun < 0 {
                    startOfLongestContiguousRun +=  numIndicesWithLowestNumIdsInRange
                }
            }
        } else {
            inContiguousRun = false
        }
    }
    // If the number of indices in range is the same for all offsets then the detected length will be wrong
    if lenOfLongestContiguousRun == numIndicesWithLowestNumIdsInRange*2 {
        lenOfLongestContiguousRun = numIndicesWithLowestNumIdsInRange
    }

    fmt.Printf(" Start of Longest Contiguous Run: %d\n", startOfLongestContiguousRun)
    fmt.Printf(" Length of Longest Contiguous Run: %d\n", lenOfLongestContiguousRun)
    fmt.Println()


    fmt.Println("Find the largest gap in the largest contiguous piece of low density")
    var startOffsetOfLowDensityRange = new(big.Int)
    startOffsetOfLowDensityRange.Mul(stepSize, big.NewInt(int64(indicesWithLowestNumIdsInRange[startOfLongestContiguousRun])))
    fmt.Printf(" Start Offset of lowest density range: %#x\n", startOffsetOfLowDensityRange)
    var rangeOfLowDensityRange = new(big.Int)
    rangeOfLowDensityRange.Mul(stepSize, big.NewInt(int64(lenOfLongestContiguousRun - 1)))
    rangeOfLowDensityRange.Add(rangeOfLowDensityRange, windowSize)
    fmt.Printf(" Size of of lowest density range: %#x\n", rangeOfLowDensityRange)

    foundIndex, foundValue := v.FindNextHighest(startOffsetOfLowDensityRange)
    fmt.Printf(" First identifier in range: Index: %d, Value: %#x\n", foundIndex, &foundValue)

    numIdsInRange := v.NumIdsInRange(startOffsetOfLowDensityRange, rangeOfLowDensityRange)
    fmt.Printf(" Number of identifiers in range: %d\n", numIdsInRange)


    var largestGap big.Int
    largestGapIndex := foundIndex
    for i := foundIndex; i < foundIndex + numIdsInRange; i++ {
        thisIndex := i % numNodes
        low := v.ids.Get(thisIndex)
        nextIndex := (i + 1) % numNodes
        // fmt.Printf("i %d, next %d\n", i, next)
        high := v.ids.Get(nextIndex)
        var diff big.Int
        if nextIndex == 0 {
            high.Add(&high, v.max)
        }
        diff.Sub(&high, &low)
        // fmt.Printf("Low:       %#x\n", &low)
        // fmt.Printf("High:      %#x\n", &high)
        // fmt.Printf("Gap %d, %d: %#x\n", this, next, &diff)
        if largestGap.Cmp(&diff) == -1 {
            largestGap = diff
            largestGapIndex = thisIndex
        }
    }

    // fmt.Printf("Location of largest gap is %d\n", largestGapIndex)
    // fmt.Printf("Largest Gap size: %#x\n", &largestGap)
    var halfOfLargestGap big.Int
    halfOfLargestGap.Div(&largestGap, big.NewInt(2))
    valStartOfLargestGapIndex := v.ids.Get(largestGapIndex)
    var ideal  big.Int
    ideal.Add(&valStartOfLargestGapIndex, &halfOfLargestGap)

    return ideal
}