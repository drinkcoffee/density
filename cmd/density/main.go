/**
 * Copyright 2020 Peter Robinson
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *   http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
    "fmt"
    "math/big"
    "github.com/drinkcoffee/density/internal/density"
)

const numNodes = 64
const dhtSpan = 16
const wordSize = 32  // 32 bytes
const numSteps = numNodes * 2


func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println()

    var gatewayIds = density.NewValueRing(wordSize)
    fmt.Println()

    fmt.Println ("Start by generating", numNodes, "random numbers")
    randomIdentifiers := []big.Int{ } 
    for i := 0; i < numNodes; i++ {
        randomIdentifiers = append(randomIdentifiers, *density.CreateRandomIdentifier())
    }
    ids := gatewayIds.GetIds()
    ids.AddMany(randomIdentifiers)
    fmt.Println(" Generated ids")
    ids.PrintAll()
    fmt.Println()



    fmt.Println("Scan the number range to see how many identifiers cover each region.")
    fmt.Println("Determine the indices that correspond to regions that have the lowest number identifiers.")
    numRange := gatewayIds.GetNumberRange()
    fmt.Printf(" Num Range: 0x%x\n", numRange)
    
    var windowSize = new(big.Int)
    windowSize.Mul(numRange, big.NewInt(dhtSpan)).Div(windowSize, big.NewInt(numNodes))
    fmt.Printf(" Window Size: 0x%x\n", windowSize)

    var stepSize = new(big.Int)
    stepSize.Div(numRange, big.NewInt(numNodes * 2))
    fmt.Printf(" Step Size: 0x%x\n", stepSize)

    fmt.Println(" Determine the weightings of each window")
    var indicesWithLowestNumIdsInRange []int = nil
    lowestNumIdsInRange := numNodes
    var offset = *big.NewInt(0)
    index := 0
    fmt.Println(" index, num identifiers in range")
    for offset.Cmp(numRange) == -1 {
        numIdsInRange := gatewayIds.NumIdsInRange(&offset, windowSize)
        // Print out comma separated values of index of the scan and the number of ids in the range.
        fmt.Printf("  %d, %d\n", index, numIdsInRange)

        if (numIdsInRange == lowestNumIdsInRange) {
            indicesWithLowestNumIdsInRange = append(indicesWithLowestNumIdsInRange, index)
        }

        if (numIdsInRange < lowestNumIdsInRange) {
            lowestNumIdsInRange = numIdsInRange
            indicesWithLowestNumIdsInRange = nil
            indicesWithLowestNumIdsInRange = append(indicesWithLowestNumIdsInRange, index)
        }

        index++
        offset.Add(&offset, stepSize)
    }

    fmt.Printf("List of window start offsets that have the lowest number of idenfiers (%d)\n", lowestNumIdsInRange)
    for _, s := range indicesWithLowestNumIdsInRange {
        fmt.Println(s)
    }
    fmt.Println()

    fmt.Println("Find the largest contiguous group of indices that have the lowest number of Ids in range")
    startOfLongestContiguousRun := 0
    lenOfLongestContiguousRun := 1
    lenOfContiguousRun := 1
    inContiguousRun := true
    numIndicesWithLowestNumIdsInRange := len(indicesWithLowestNumIdsInRange)
    // Go through the array twice to ensure we handle the wrap-around
    for i := 1; i < numIndicesWithLowestNumIdsInRange * 2; i++ {
        idx := i % numIndicesWithLowestNumIdsInRange
        idxLessOne := (i-1) % numIndicesWithLowestNumIdsInRange
        if (indicesWithLowestNumIdsInRange[idx] - indicesWithLowestNumIdsInRange[idxLessOne] == 1) ||
           (indicesWithLowestNumIdsInRange[idxLessOne] == numSteps-1 && indicesWithLowestNumIdsInRange[idx] == 0) {
            // Is contiguous
            if (inContiguousRun) {
                lenOfContiguousRun++
                if (lenOfLongestContiguousRun < lenOfContiguousRun) {
                    lenOfLongestContiguousRun = lenOfContiguousRun
                }
            } else {
                lenOfContiguousRun = 1
                inContiguousRun = true
            }
            startOfLongestContiguousRun = (idx - lenOfContiguousRun) % numIndicesWithLowestNumIdsInRange
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
    fmt.Printf(" Start Offset of lowest density range: 0x%x\n", startOffsetOfLowDensityRange)
    var rangeOfLowDensityRange = new(big.Int)
    rangeOfLowDensityRange.Mul(stepSize, big.NewInt(int64(lenOfLongestContiguousRun)))
    rangeOfLowDensityRange.Add(rangeOfLowDensityRange, windowSize)
    fmt.Printf(" Size of of lowest density range: 0x%x\n", rangeOfLowDensityRange)

    foundIndex, foundValue := gatewayIds.FindNextHighest(startOffsetOfLowDensityRange)
    fmt.Printf(" First identifier in range: Index: %d, Value: 0x%x\n", foundIndex, &foundValue)

    numIdsInRange := gatewayIds.NumIdsInRange(startOffsetOfLowDensityRange, rangeOfLowDensityRange)
    fmt.Printf(" Number of identifiers in range: %d\n", numIdsInRange)


    var largestGap big.Int
    largestGapIndex := foundIndex
    for i := foundIndex; i < foundIndex + numIdsInRange; i++ {
        this := i % numNodes
        low := ids.Get(this)
        next := (i + 1) % numNodes
        fmt.Printf("i %d, next %d\n", i, next)
        high := ids.Get(next)
        var diff big.Int
        diff.Sub(&high, &low)
        fmt.Printf("Low:       0x%x\n", &low)
        fmt.Printf("High:      0x%x\n", &high)
        fmt.Printf("Diff %d, %d: 0x%x\n", i, i+1, &diff)

        if largestGap.Cmp(&diff) == -1 {
            largestGap = diff
            largestGapIndex = i
        }
    }

    var halfOfLargestGap big.Int
    halfOfLargestGap.Div(&largestGap, big.NewInt(2))
    valStartOfLargestGapIndex := ids.Get(largestGapIndex)
    var ideal  big.Int
    ideal.Add(&valStartOfLargestGapIndex, &halfOfLargestGap)


    fmt.Printf("Ideal:     0x%x\n", &ideal)

    




}