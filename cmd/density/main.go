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



func main() {
    fmt.Println("Gateway Id Density")
    fmt.Println()

    fmt.Println ("Start by generating", numNodes, "random numbers")
    var gatewayIds = density.NewGatewayIds(wordSize)
    for i := 0; i < numNodes; i++ {
        gatewayIds.AddRandom()
    }
    fmt.Println()

    fmt.Println ("Sort them")
    gatewayIds.Sort();
    gatewayIds.PrintAll()


    var numRange = new(big.Int)
    numRange.Exp(big.NewInt(2), big.NewInt(8 * wordSize), nil)
    fmt.Printf("Num Range: 0x%x\n", numRange)
    
    var windowSize = new(big.Int)
    windowSize.Mul(numRange, big.NewInt(dhtSpan)).Div(windowSize, big.NewInt(numNodes))
    fmt.Printf("Window Size: 0x%x\n", windowSize)

    var stepSize = new(big.Int)
    stepSize.Div(numRange, big.NewInt(numNodes * 2))
    fmt.Printf("Step Size: 0x%x\n", stepSize)

    fmt.Println("Weightings")
    var indicesWithLowestNumIdsInRange []int = nil
    lowestNumIdsInRange := numNodes
    var offset = *big.NewInt(0)
    index := 0
    for offset.Cmp(numRange) == -1 {
        numIdsInRange := gatewayIds.NumIdsInRange(&offset, windowSize)
        // Print out comma separated values of index of the scan and the number of ids in the range.
        fmt.Printf("%d, %d\n", index, numIdsInRange)

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

    for i, s := range indicesWithLowestNumIdsInRange {
        fmt.Println(i, s)
    }

    // Find the largest contiguous group of indices that have the lowest number of Ids in range
    startOfLongestContiguousRun := 0
    longestLenOfContiguousRun := 1
    lenOfContiguousRun := 1
    inContiguousRun := true
    for i := 1; i < len(indicesWithLowestNumIdsInRange); i++ {
        if (indicesWithLowestNumIdsInRange[i] - indicesWithLowestNumIdsInRange[i-1] == 1) {
            // Is contiguous
            if (inContiguousRun) {
                lenOfContiguousRun++
                if (longestLenOfContiguousRun < lenOfContiguousRun) {
                    longestLenOfContiguousRun = lenOfContiguousRun
                    startOfLongestContiguousRun = i - lenOfContiguousRun
                }
            } else {
                lenOfContiguousRun = 1
                inContiguousRun = true
            }
        } else {
            inContiguousRun = false
        }
    }
    fmt.Printf("startOfLongestContiguousRun: %d\n", startOfLongestContiguousRun)
    fmt.Printf("longestLenOfContiguousRun: %d\n", longestLenOfContiguousRun)

    // Find the largest gap in the contiguous run of low density
    var largestGap big.Int
    largestGapIndex := startOfLongestContiguousRun
    for i := startOfLongestContiguousRun; i < startOfLongestContiguousRun + longestLenOfContiguousRun; i++ {
        low := gatewayIds.Get(i)
        next := (i + 1) % numNodes
        fmt.Printf("i %d, next %d\n", i, next)
        high := gatewayIds.Get(next)
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
    valStartOfLargestGapIndex := gatewayIds.Get(largestGapIndex)
    var ideal  big.Int
    ideal.Add(&valStartOfLargestGapIndex, &halfOfLargestGap)


    fmt.Printf("Ideal:     0x%x\n", &ideal)

    




}